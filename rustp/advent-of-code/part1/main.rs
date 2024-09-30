use std::{fs::File, io::{BufRead, BufReader}};




fn main() {
    let reader: BufReader<File> = BufReader::new(File::open("1-input.txt").expect("Cannot open file"));
    let mut sum: i32 = 0;
    for line in reader.lines() {
        let result: Vec<char> = find_numbers(line.unwrap().as_str());

        let result_string: String;
        if result.is_empty() {
            return;
        } else if result.len() < 2 {
            result_string = format!("{}{}", result.first().unwrap(), result.first().unwrap());
        } else {
            result_string = format!("{}{}", result.first().unwrap(), result.last().unwrap());
        }

        let my_int: i32 =  result_string.parse::<i32>().unwrap();
        sum += my_int;
    }

    println!("Sum: {}", sum);
}

// receive a string and return a string of characters
fn find_numbers(s: &str) -> Vec<char> {
    let mut result = Vec::new();
    for c in s.chars() {
        if c.is_numeric() {
            result.push(c);
        }
    }
    result
}
