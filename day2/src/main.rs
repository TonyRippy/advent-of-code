use std::fs::File;
use std::io::BufRead;


fn parse_line(line: &str) -> (usize, (usize, usize, usize)) {
    let (mut red, mut green, mut blue) = (0,0,0);

    assert_eq!(&line[0..5], "Game ");
    let i = line.find(':').unwrap();
    let game = line[5..i].parse::<usize>().unwrap();
    for set in line[i+1..].split(';') {
        for sample in set.split(',').map(str::trim) {
            let tokens = sample.split_ascii_whitespace().collect::<Vec<_>>();
            assert_eq!(tokens.len(), 2);
            let n = tokens[0].parse::<usize>().unwrap();
            match tokens[1] {
                "red" => {
                    if n > red {
                        red = n;
                    }
                },
                "green" => {
                    if n > green {
                        green = n;
                    }
                },
                "blue" => {
                    if n > blue {
                        blue = n;
                    }
                },
                _ => panic!("Unknown color: {}", tokens[1]),
            }
        }
    }
    (game, (red, green, blue))
}

fn is_possible(guess: (usize, usize, usize), game: (usize, usize, usize)) -> bool {
    game.0 <= guess.0 && game.1 <= guess.1 && game.2 <= guess.2
}

fn power(game: (usize, usize, usize)) -> usize {
    game.0 * game.1 * game.2
}

fn main() -> Result<(), std::io::Error> {
    // let fname = std::env::args().nth(1).unwrap();
    let fname = "input.txt";

    let guess = (12,13,14);
    let mut sum = 0;

    let file = File::open(fname)?;
    let reader = std::io::BufReader::new(file);
    for l in reader.lines() {
        let line = l.unwrap();
        let (id, game) = parse_line(&line);
        sum += power(game);
    }
    println!("SUM = {}", sum);
    Ok(())
}
