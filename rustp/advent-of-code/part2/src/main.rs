use std::{collections::HashMap,  fs::File, io::{BufRead, BufReader}};

fn main() {
    let _dict: HashMap<&str, &str> = HashMap::from(
        [
            ("one", "1"),
            ("two", "2"),
            ("three", "3"),
            ("four", "4"),
            ("five", "5"),
            ("six", "6"),
            ("seven", "7"),
            ("eight", "8"),
            ("nine", "9"),
        ]
    );

    let mut sum: i32 = 0;
    let reader: BufReader<File> = BufReader::new(File::open("./src/input-real.txt").expect("Cannot open file"));
    for line in reader.lines() {
        println!("Processing line: {:?}", line);
        let line: String = line.unwrap();
        let mut _dict2: HashMap<usize, char> = HashMap::new();
        let result_int = fun_name2(&_dict, line, _dict2);
        sum += result_int;
        println!("Sum: {}", sum);
    }
    println!("Sum: {}", sum);
}

fn fun_name2(_dict: &HashMap<&str, &str>, line: String, mut _dict2: HashMap<usize, char>) -> i32 {
    fun_name(_dict, &line, &mut _dict2);

    fun_name1(line, &mut _dict2);

    println!("{:?}", _dict2);

    let max_index: usize = *_dict2.keys().max().unwrap();
    let min_index: usize = *_dict2.keys().min().unwrap();

    let result_string: String;

    result_string = format!("{}{}", _dict2.get(&min_index).unwrap(), _dict2.get(&max_index).unwrap());
    println!("Result string: {}", result_string);
    let result_int: i32 = result_string.parse::<i32>().unwrap();
    result_int
}

fn fun_name1(line: String, _dict2: &mut HashMap<usize, char>) {
    for (index, value) in line.chars().enumerate() {
        if value.is_numeric() {
            _dict2.insert(index, value);
        }
    }
}

fn fun_name(number_string_2_number_int: &HashMap<&str, &str>, line: &String, _index_2_number: &mut HashMap<usize, char>) {
    for (key, _value) in number_string_2_number_int.iter() {
        let matches: Vec<(usize, &str)> = line.match_indices(key).collect();
        if matches.len() > 0 {
            for (index, _vv) in matches.iter() {
                // as value just have one chart we can use next unwrap without any problem
                _index_2_number.insert(*index, _value.chars().next().unwrap());
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::collections::HashMap;

    #[test]
    fn test_fun_name2() {
        let mut _dict2: HashMap<usize, char> = HashMap::new();
        let _dict: HashMap<&str, &str> = HashMap::from(
            [
                ("one", "1"),
                ("two", "2"),
                ("three", "3"),
                ("four", "4"),
                ("five", "5"),
                ("six", "6"),
                ("seven", "7"),
                ("eight", "8"),
                ("nine", "9"),
            ]
        );
        let line: String = "threerznlrhtkjp23mtflmbrzq395three".to_string();
        let result_int: i32 = fun_name2(&_dict, line, _dict2);
        assert_eq!(result_int, 33);
    }

    #[test]
    fn test_fun_name1() {
        let _dict: HashMap<&str, &str> = HashMap::from(
                    [
                        ("one", "1"),
                        ("two", "2"),
                        ("three", "3"),
                        ("four", "4"),
                        ("five", "5"),
                        ("six", "6"),
                        ("seven", "7"),
                        ("eight", "8"),
                        ("nine", "9"),
                    ]
                );
        let mut _dict2: HashMap<usize, char> = HashMap::new();
        let line: String = "twotwo".to_string();
        fun_name(&_dict, &line, &mut _dict2);
        assert_eq!(_dict2.get(&0).unwrap(), &'2');
        assert_eq!(_dict2.get(&3).unwrap(), &'2');

        let line: String = "twotwothree".to_string();
        fun_name(&_dict, &line, &mut _dict2);

        assert_eq!(_dict2.get(&0).unwrap(), &'2');
        assert_eq!(_dict2.get(&3).unwrap(), &'2');
        assert_eq!(_dict2.get(&6).unwrap(), &'3');
    }
}
