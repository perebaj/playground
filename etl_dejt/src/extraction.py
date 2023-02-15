import glob
import json
import re
import tempfile

import fitz


class Extraction():
    """Class for extracting text from pdf files."""
    def __init__(self, file_path):
        self.file_path = file_path
        self.file_name = file_path.split("/")[-1]
        self.output_path = f"data/{self.file_name}_output.txt"
        self._PROCESS_IDENTIFIER = "PROCESSO.*\d{7}-\d{2}.\d{4}.\d{1}.\d{2}.\d{4}"
        self._REMOVE_REGEX_1 = "^Código para aferir autenticidade deste caderno: \d{6}$"
        self._REMOVE_REGEX_2 = "^\d{4}\/\d{4}$"
        self._REMOVE_REGEX_3 = "^\d+$"
        self._REMOVE_REGEX_4 = "^Tribunal Superior do Trabalho$"

    def pdf_to_text(self):
        """Extract text from pdf file."""
        doc = fitz.open(self.file_path)
        temp_file = tempfile.NamedTemporaryFile()
        text: str = ""
        for page in doc:
            text = page.get_text()
            temp_file.write(text.encode("utf-8"))
        temp_file.seek(0)
        with open(self.output_path, "wb") as f:
            f.write(temp_file.read())
            # print(f"Saved output:{self.output_path} ")
        temp_file.close()

    def get_output_txt(self) -> list[str]:
        process_text = open(self.output_path, "r")
        process_list = process_text.readlines()
        return process_list

    def process_chunk(self):
        process_list = self.get_output_txt()
        process_chunk_list = []
        process_chunk: str = ""
        for line in process_list:
            if re.match(self._PROCESS_IDENTIFIER, line, re.IGNORECASE):
                # process_chunk = process_chunk + line
                process_chunk_list.append(process_chunk)
                process_chunk = "" + line
            else:
                process_chunk = process_chunk + line
        return process_chunk_list

    def clean_chunk(self, process_chunk_list):
        clean_chunk_list = []
        for chunk in process_chunk_list:
            new_chunk = re.sub(self._REMOVE_REGEX_1, "", chunk, flags=re.MULTILINE)
            new_chunk = re.sub(self._REMOVE_REGEX_2, "", new_chunk, flags=re.MULTILINE)
            new_chunk = re.sub(self._REMOVE_REGEX_3, "", new_chunk, flags=re.MULTILINE)
            new_chunk = re.sub(self._REMOVE_REGEX_4, "", new_chunk, flags=re.MULTILINE)
            clean_chunk_list.append(new_chunk)
        return clean_chunk_list

    def mount_table(self, process_chunk_list):
        table_list = []
        for process in process_chunk_list:
            dict_table = {}
            match = re.search(self._PROCESS_IDENTIFIER, process, re.IGNORECASE)
            if match:
                dict_table["procces_id"] = match.group()
            else:
                dict_table["procces_id"] = None
            dict_table["process"] = process
            dict_table["extract_at"] = self.get_extracted_at()
            dict_table["title"] = self.file_name
            dict_table["book"] = "Caderno do Tribunal Superior do Trabalho - Judiciário"
            
            if len(process) > 2000:
                table_list.append(dict_table)
        return table_list

    def get_extracted_at(self):
        match = re.search(r"Diario_(\d+)__(\d+)_(\d+)_(\d+)", self.file_name)
        if match:
            day = match.group(2)
            month = match.group(3)
            year = match.group(4)
            return f"{year}-{month}-{day}"
        else:
            return None


if __name__ == "__main__":
    path = "/root/playground/etl_dejt/data/data17-01-2023--27-01-2023/*.pdf"
    files_list = glob.glob(path)
    entiry_table_list = []
    for file in files_list:
        extraction = Extraction(file)
        extraction.pdf_to_text()
        process_chunk_list = extraction.process_chunk()
        print(f"Process chunk list size: {len(process_chunk_list)}")
        clean_chunk_list = extraction.clean_chunk(process_chunk_list)
        table_list = extraction.mount_table(clean_chunk_list)
        print(f"Table list size: {len(table_list)}")
        entiry_table_list.extend(table_list)
 

    clean_path = f"clean/entire_table_clean3.json"
    with open(clean_path, "w") as f:
        json.dump(entiry_table_list, f, ensure_ascii=False)
