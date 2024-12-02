/*

Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green
Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue
Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red
Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red
Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green

only 12 red cubes, 13 green cubes, and 14 blue cubes

1) interpretar os dados
2) fazer um for em cada set de elementos.
3) Se no set tiver elementos de mesma cor, somar a quantidade de elementos
4) Comparação com meu map de cores

*/

use std::collections::HashMap;
use regex::Regex;

fn main() {
    // let game: &str = "Game 1: 3 blue, 4 red, 1 blue; 1 red, 2 green, 6 blue; 2 green";

    let games: Vec<&str> = vec![
        "Game 1: 3 blue, 4 red; 1 red, 2 green, 30 blue; 2 green",
        // "Game 200: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue",
        // "Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red",
        // "Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red",
        // "Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green",
    ];

    for game in games {
        let mut sum: i32 = 0;
        // The first position will fit the Game x, and the second the rest of the string
        let itera: Vec<&str> = game.split(":").collect::<Vec<&str>>();
        let re: Regex = Regex::new(r"\d+").unwrap();
        let caps: regex::Captures<'_> = re.captures(itera[0]).unwrap();
        let game_index: &str = caps.get(0).unwrap().as_str();

        let games: Vec<&str> = itera[1].split(";").collect::<Vec<&str>>();

        let _dict: HashMap<&str, i32> = HashMap::from([("red", 12), ("green", 13), ("blue", 14)]);

        for i in games {
            let sub_games: Vec<&str> = i.trim().split(",").collect();
            let mut _dict2: HashMap<&str, i32> = HashMap::new();
            for game in sub_games {
                let re: Regex = Regex::new(r"(\d+) (\w+)").unwrap();
                let color_number: regex::Captures<'_> = re.captures(game).unwrap();
                let number: &str = color_number.get(1).unwrap().as_str();
                let color: &str = color_number.get(2).unwrap().as_str();
                _dict2.insert(color, number.parse::<i32>().unwrap());
            }

            println!("{:?}", _dict2);

            for (key, vv) in _dict.iter() {
                let aux_val: Option<&i32> = _dict2.get(key);
                if aux_val.is_some() && aux_val.unwrap() > vv {
                    println!("invalid game: {}", game_index);
                } else {
                    sum = sum + game_index.parse::<i32>().unwrap()
                }
            }

            println!("sum {}", sum)
        }
    }
}
