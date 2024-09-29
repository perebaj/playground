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
        for (key, value) in _dict.iter() {
            let string_index: Option<usize> = line.find(key);
            if string_index.is_some() {
                _dict2.insert(string_index.unwrap(), value.chars().next().unwrap());
            }
        }

        for (index, value) in line.chars().enumerate() {
            if value.is_numeric() {

                _dict2.insert(index, value);
            }
        }

        println!("{:?}", _dict2);

        let max_index: usize = *_dict2.keys().max().unwrap();
        let min_index: usize = *_dict2.keys().min().unwrap();

        println!("Max index: {}", max_index);
        println!("Min index: {}", min_index);

        let result_string: String;

        result_string = format!("{}{}", _dict2.get(&min_index).unwrap(), _dict2.get(&max_index).unwrap());
        println!("Result string: {}", result_string);
        let result_int: i32 = result_string.parse::<i32>().unwrap();
        sum += result_int;
        println!("Result: {}", sum);
    }
    println!("Sum: {}", sum);
}

/*

eighttwothree

823


*/
