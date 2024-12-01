use std::fs::File;
use std::io::BufRead;

const EMPTY_COST: usize = 1000000 - 1;

fn read_pairs(reader: impl BufRead) -> Result<Vec<(usize, usize)>, std::io::Error> {
    let mut v = Vec::new();
    let mut y = 0usize;
    let mut col_seen = Vec::new();
    for line in reader.lines() {
        let l = line?;
        if col_seen.is_empty() {
            col_seen.resize(l.len(), false);
        }
        let seen = l
            .chars()
            .enumerate()
            .filter(|(_, c)| *c == '#')
            .map(|(i, _)| i)
            .collect::<Vec<_>>();
        if seen.is_empty() {
            y += EMPTY_COST;
        } else {
            for x in seen.into_iter() {
                v.push((x, y));
                col_seen[x] = true;
            }
        }
        y += 1;
    }
    // dbg!(&col_seen);
    let mut offsets = Vec::new();
    let mut offset = 0usize;
    for seen in col_seen.into_iter() {
        offsets.push(offset);
        offset += if seen { 0 } else { EMPTY_COST };
    }
    // dbg!(&offsets);
    Ok(v.into_iter().map(|(x, y)| (x + offsets[x], y)).collect())
}

fn main() -> Result<(), std::io::Error> {
    let fname = std::env::args().nth(1).unwrap();
    let file = File::open(fname)?;
    let reader = std::io::BufReader::new(file);

    let pairs = read_pairs(reader)?;
    // dbg!(pairs);

    let mut sum = 0;
    for (i, p1) in pairs.iter().enumerate() {
        for (j, p2) in pairs.iter().enumerate() {
            if j < i {
                continue;
            }
            let (x1, y1) = *p1;
            let (x2, y2) = *p2;
            let dx = x2 as i64 - x1 as i64;
            let dy = y2 as i64 - y1 as i64;
            let dist = dx.abs() + dy.abs();
            // println!("{}, {} -> {:?}, {:?} --> {}", i, j, p1, p2, dist);
            sum += dist;
        }
    }
    println!("sum: {}", sum);

    // TODO: Solve a problem!
    Ok(())
}
