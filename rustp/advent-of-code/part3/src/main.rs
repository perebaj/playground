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

fn main() {
    // let game: &str = "Game 1: 3 blue, 4 red, 1 blue; 1 red, 2 green, 6 blue; 2 green";

    let games = vec![
        "Game 1: 3 blue, 4 red; 1 red, 2 green, 30 blue; 2 green",
        // "Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue",
        // "Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red",
        // "Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red",
        // "Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green",
    ];

    for game in games {
        let mut sum: i32 = 0;
        let itera: Vec<&str> = game.split(":").collect::<Vec<&str>>();
        let game_index: &str = &itera[0][5..6];
        let itera2: Vec<&str> = itera[1].split(";").collect::<Vec<&str>>();

        let _dict: HashMap<&str, i32> = HashMap::from([("red", 12), ("green", 13), ("blue", 14)]);

        for i in itera2 {
            let itera3: Vec<&str> = i.trim().split(",").collect();
            let mut _dict2: HashMap<&str, i32> = HashMap::new();
            for i in itera3 {
                let i_aux: &str = i.trim();

                let val: i32 = i_aux.split(" ").collect::<Vec<&str>>()[0].parse::<i32>().unwrap();
                let key: &str = &i_aux[2..(i_aux.len())].trim();
                if let Some(val_aux) = _dict2.get(key) {
                    _dict2.insert(key, val_aux + val);
                } else {
                    _dict2.insert(key, val);
                }
            }
            println!("{:?}", _dict2);
            // iterate over _dict and compare with _dict2
            for (key, vv) in _dict.iter() {
                let aux_val: Option<&i32> = _dict2.get(key);
                println!("Key: {}, Value: {}, aux_var {:?}", key, vv, aux_val);
                if aux_val.is_some() && aux_val.unwrap() > vv {
                    println!("invalid game: {}", game_index);
                }
            }

            println!("Valid game: {}", game_index);
            println!("sum : {}", sum)
        }
    }
}
