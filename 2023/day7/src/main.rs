use std::cmp::Ordering;
use std::fs::File;
use std::io::BufRead;

fn value_of_card(name: char) -> usize {
    match name {
        'A' => 14,
        'K' => 13,
        'Q' => 12,
        'J' => 1, // was 11 in part 1
        'T' => 10,
        _ => name.to_digit(10).unwrap() as usize,
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord)]
enum HandType {
    HighCard = 1,
    OnePair = 2,
    TwoPairs = 3,
    ThreeOfAKind = 4,
    FullHouse = 5,
    FourOfAKind = 6,
    FiveOfAKind = 7,
}

#[derive(Debug)]
struct Hand {
    name: String,
    cards: Vec<usize>,
    bid: usize,
    hand_type: HandType,
}

impl Hand {
    pub fn parse(line: &str) -> Self {
        let mut parts = line.split_whitespace();
        let name = parts.next().unwrap().to_string();
        let cards = name.chars().map(value_of_card).collect::<Vec<usize>>();
        let bid = parts.next().unwrap().parse::<usize>().unwrap();
        let hand_type = Self::hand_type(&cards);
        Self {
            name,
            cards,
            bid,
            hand_type,
        }
    }

    pub fn hand_type(cards: &[usize]) -> HandType {
        let mut groups = vec![];
        let jokers = cards.iter().filter(|x| **x == 1).count();
        if jokers == 5 {
            return HandType::FiveOfAKind;
        }
        let mut cs = cards
            .iter()
            .filter_map(|&x| match x {
                1 => None,
                _ => Some(x),
            })
            .collect::<Vec<_>>();
        cs.sort();
        let mut i = 0;
        while i < cs.len() {
            let v = cs[i];
            let mut count = 1;
            i += 1;
            while i < cs.len() && cs[i] == v {
                count += 1;
                i += 1;
            }
            groups.push((v, count));
        }
        groups.sort_by(|a, b| {
            let cmp = b.1.cmp(&a.1);
            if cmp != Ordering::Equal {
                return cmp;
            }
            b.0.cmp(&a.0)
        });
        groups[0].1 += jokers;
        match groups.len() {
            1 => HandType::FiveOfAKind,
            2 => match groups[0].1 {
                4 => HandType::FourOfAKind,
                3 => HandType::FullHouse,
                _ => panic!("Invalid hand!"),
            },
            3 => match groups[0].1 {
                3 => HandType::ThreeOfAKind,
                2 => HandType::TwoPairs,
                _ => panic!("Invalid hand!"),
            },
            4 => HandType::OnePair,
            5 => HandType::HighCard,
            _ => panic!("Invalid hand!"),
        }
    }

    fn cmp(&self, other: &Self) -> Ordering {
        let cmp = self.hand_type.cmp(&other.hand_type);
        if cmp != Ordering::Equal {
            return cmp;
        }
        for (s, o) in self.cards.iter().zip(other.cards.iter()) {
            let cmp = s.cmp(o);
            if cmp != Ordering::Equal {
                return cmp;
            }
        }
        Ordering::Equal
    }
}

fn main() -> Result<(), std::io::Error> {
    let fname = std::env::args().nth(1).unwrap();
    let file = File::open(fname)?;
    let reader = std::io::BufReader::new(file);

    let mut hands = reader
        .lines()
        .map(|x| Hand::parse(x.unwrap().as_str()))
        .collect::<Vec<_>>();
    hands.sort_by(Hand::cmp);

    let winnings: usize = hands
        .iter()
        .enumerate()
        .map(|(rank, hand)| {
            let score = (rank + 1) * hand.bid;
            println!("Rank: {} {:?}, Score: {}", rank + 1, hand, score);
            score
        })
        .sum();

    println!("Total winnings: {}", winnings);
    // TODO: Solve a problem!
    Ok(())
}

mod test {
    use super::*;

    #[test]
    fn test_one_pair() {
        let input = vec![13, 10, 3, 3, 2];
        let hand_type = Hand::hand_type(&input);
        assert_eq!(hand_type, HandType::OnePair);
    }

    #[test]
    fn test_joker_with_pair() {
        let input = vec![13, 1, 3, 3, 2];
        let hand_type = Hand::hand_type(&input);
        assert_eq!(hand_type, HandType::ThreeOfAKind);
    }

    #[test]
    fn test_joker_with_high_card() {
        let input = vec![13, 1, 4, 3, 2];
        let hand_type = Hand::hand_type(&input);
        assert_eq!(hand_type, HandType::OnePair);
    }

    #[test]
    fn test_joker_with_two_pairs() {
        let input = vec![2, 1, 4, 2, 4];
        let hand_type = Hand::hand_type(&input);
        assert_eq!(hand_type, HandType::FullHouse);
    }
}
