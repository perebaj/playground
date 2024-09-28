import itertools

import pandas as pd
from optbinning import OptimalBinning, OptimalBinning2D
from sklearn.base import BaseEstimator, TransformerMixin
from sklearn.utils.validation import check_is_fitted


class WoeCorrelationTransformer(BaseEstimator, TransformerMixin):
    def __init__(
        self,
        gamma=0,
        iv_filter=False,
        iv_threshold=0.02,
        monotonic_trend="auto_asc_desc",
        class_weight="balanced",
        combine=False,
        corr_method="spearman",
        corr_threshold=0.5,
        to_drop_features=[],
        selected_features=[],
    ):
        if not isinstance(gamma, (int, float)):
            raise TypeError("gamma must be an int or float instance")

        if not isinstance(iv_filter, bool):
            raise TypeError("iv_filter must be a bool instance")

        if not isinstance(iv_threshold, float):
            raise TypeError("iv_threshold must be a float instance")

        if not isinstance(monotonic_trend, str):
            raise TypeError("monotonic_trend must be a str instance")

        if not isinstance(class_weight, str) and (class_weight is not None):
            raise TypeError("class_weight must be a str instance or None.")

        if not isinstance(combine, bool):
            raise TypeError("combine must be a bool instance")

        if not isinstance(corr_method, str):
            raise TypeError("corr_method must be a str instance")

        if not isinstance(corr_threshold, float):
            raise TypeError("corr_threshold must be a float instance")

        if not isinstance(to_drop_features, list):
            raise TypeError("to_drop_features must be a list instance")

        if not isinstance(selected_features, list):
            raise TypeError("selected_features must be a list instance")

        self.gamma = gamma
        self.iv_filter = iv_filter
        self.iv_threshold = iv_threshold
        self.monotonic_trend = monotonic_trend
        self.class_weight = class_weight
        self.combine = combine
        self.corr_method = corr_method
        self.corr_threshold = corr_threshold
        self.combined_feature_splitter = "_|X|_"
        self.to_drop_features = to_drop_features
        self.selected_features = selected_features

        self.woe_encoder_binners = None
        self.woe_combiner_binners = None
        self.corr_filtered_features = None

        self._is_woe_encoder_fitted = False
        self._is_woe_combiner_fitted = False
        self._is_corr_filter_fitted = False

        self._is_fitted = False

    def __sklearn_is_fitted__(self):
        return self._is_fitted

    def fit_woe_encoder(self, X, y, sample_weight=None):
        _X = X.copy()

        numerical_features = _X.select_dtypes("number").columns.tolist()

        woe_encoder_binners = {}

        for feature in _X.columns:
            if feature in self.to_drop_features:
                continue

            if len(_X[feature].unique()) < 2:
                continue

            feature_dtype = "numerical" if feature in numerical_features else "categorical"

            binner = OptimalBinning(
                dtype=feature_dtype,
                solver="mip",
                gamma=self.gamma,
                divergence="iv",
                monotonic_trend=self.monotonic_trend,
                min_prebin_size=0.0001,
                class_weight=self.class_weight,
                random_state=42,
            )
            binner = binner.fit(
                _X[feature].copy(),
                y.copy(),
                sample_weight=(sample_weight.copy() if sample_weight is not None else None),
            )
            binning_table = binner.binning_table.build()

            feature_iv = binning_table["IV"]["Totals"]

            if self.iv_filter:
                if feature_iv > self.iv_threshold:
                    woe_encoder_binners.update({feature: binner})

            else:
                woe_encoder_binners.update({feature: binner})

        self.woe_encoder_binners = woe_encoder_binners
        self._is_woe_encoder_fitted = True

        return self

    def fit_woe_combiner(self, X, y):
        _X = X.copy()

        numerical_features = _X.select_dtypes("number").columns.tolist()

        woe_combiner_binners = {}

        for feature_1, feature_2 in itertools.combinations(sorted(list(self.woe_encoder_binners.keys())), 2):
            combined_feature_name = f"COMB_{feature_1}{self.combined_feature_splitter}{feature_2}"

            if combined_feature_name in self.to_drop_features:
                continue

            if (combined_feature_name not in self.selected_features) and (len(self.selected_features) > 0):
                continue

            feature_1_dtype = "numerical" if feature_1 in numerical_features else "categorical"
            feature_2_dtype = "numerical" if feature_2 in numerical_features else "categorical"

            binner = OptimalBinning2D(
                dtype_x=feature_1_dtype,
                dtype_y=feature_2_dtype,
                solver="mip",
                gamma=self.gamma,
                divergence="iv",
                min_prebin_size_x=0.0001,
                min_prebin_size_y=0.0001,
            )

            try:
                binner = binner.fit(_X[feature_1], _X[feature_2], y)

            except (IndexError, ValueError):
                print(f"Problem combining {feature_1} with {feature_2}")

                continue

            binning_table = binner.binning_table.build()

            feature_1_iv = self.woe_encoder_binners[feature_1].binning_table.build()["IV"]["Totals"]
            feature_2_iv = self.woe_encoder_binners[feature_2].binning_table.build()["IV"]["Totals"]

            combined_feature_iv = binning_table["IV"]["Totals"]

            if combined_feature_iv > (feature_1_iv + feature_2_iv):
                woe_combiner_binners.update({combined_feature_name: binner})

            elif combined_feature_name in self.selected_features:
                woe_combiner_binners.update({combined_feature_name: binner})

        self.woe_combiner_binners = woe_combiner_binners
        self._is_woe_combiner_fitted = True

        return self

    def unfitted_transform(self, X, y=None):
        _X = X.copy()

        for (
            combined_feature,
            combined_binner,
        ) in self.woe_combiner_binners.items():
            feature_1, feature_2 = combined_feature.replace("COMB_", "").split(self.combined_feature_splitter)

            try:
                _X[combined_feature] = combined_binner.transform(_X[feature_1], _X[feature_2], metric="woe")

            except Exception:
                raise SystemError(f"Problem with {combined_feature} during encoding")

        necessary_features = list(self.woe_encoder_binners.keys()) + list(self.woe_combiner_binners.keys())
        _X = _X[necessary_features]

        for feature, individual_binner in self.woe_encoder_binners.items():
            try:
                _X[feature] = individual_binner.transform(_X[feature], metric="woe")

            except Exception:
                raise SystemError(f"Problem with {feature} during encoding")

        return _X

    def fit_corr_filter(self, X, y=None):
        _X = X.copy()

        _X = self.unfitted_transform(_X)

        corr_df = _X.corr(method=self.corr_method)
        corr_clusters_dict = {}
        corr_cluster_index = 0

        while not corr_df.empty:
            selected_feature = corr_df.columns.tolist()[0]

            corr_series = corr_df[selected_feature].abs().copy()
            filtered_corr_series = corr_series[corr_series >= self.corr_threshold].drop(selected_feature)
            correlated_features = filtered_corr_series.index.tolist() + [selected_feature]

            corr_clusters_dict.update({f"cluster_{corr_cluster_index}": correlated_features})
            corr_cluster_index += 1
            corr_df = corr_df.drop(index=correlated_features, columns=correlated_features)

        corr_filtered_features = []
        all_binners_dict = {
            **self.woe_encoder_binners,
            **self.woe_combiner_binners,
        }

        for clustered_features in corr_clusters_dict.values():
            if len(clustered_features) > 1:
                feature_iv_series = pd.Series(
                    [all_binners_dict[feature].binning_table.build()["IV"]["Totals"] for feature in clustered_features]
                )

                filtered_feature_index = feature_iv_series.idxmax()
                filtered_feature = clustered_features[filtered_feature_index]

            else:
                filtered_feature = clustered_features[0]

            corr_filtered_features.append(filtered_feature)

        self.corr_filtered_features = corr_filtered_features
        self._is_corr_filter_fitted = True

        return self

    def fit_manual_feature_selection(self):
        filtered_woe_encoder_binners = {
            feature: binner for feature, binner in self.woe_encoder_binners.items() if feature in self.selected_features
        }
        filtered_woe_combiner_binners = {
            feature: binner
            for feature, binner in self.woe_combiner_binners.items()
            if feature in self.selected_features
        }

        if self.corr_filtered_features is not None:
            filtered_corr_filtered_features = [
                feature for feature in self.corr_filtered_features if feature in self.selected_features
            ]

        else:
            filtered_corr_filtered_features = list(filtered_woe_encoder_binners.keys()) + list(
                filtered_woe_combiner_binners.keys()
            )

        self.woe_encoder_binners = filtered_woe_encoder_binners
        self.woe_combiner_binners = filtered_woe_combiner_binners
        self.corr_filtered_features = filtered_corr_filtered_features

        return self

    def fit(self, X, y, sample_weight=None):
        _X = X.copy()

        self.fit_woe_encoder(_X, y, sample_weight)

        if self.combine:
            self.fit_woe_combiner(_X, y)

        else:
            self.woe_combiner_binners = {}

        if self.corr_threshold < 1.0:
            self.fit_corr_filter(_X, y)

        if len(self.selected_features) > 0:
            self.fit_manual_feature_selection()

        condition_1 = self._is_woe_encoder_fitted
        condition_2 = self.combine == self._is_woe_combiner_fitted
        condition_3 = (self.corr_threshold < 1.0) == self._is_corr_filter_fitted

        self._is_fitted = condition_1 and condition_2 and condition_3

        return self

    def transform(self, X, y=None):
        check_is_fitted(self)

        _X = X.copy()

        _X = self.unfitted_transform(_X)
        _X = _X[self.corr_filtered_features] if self.corr_filtered_features is not None else _X

        # Solving some bug with OptimalBinning with unknown values
        _X = _X.fillna(0.0)

        return _X
