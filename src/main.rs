use axum::Router;

struct AppState {
    next_user_id: u32,
    lobbies: Vec<Lobby>,
}

fn main() {
    let game_data = game_data::GameData::default();

    for province in &game_data.provinces {
        println!("{}", province.name);
        let mut count = 0;
        for adj in &game_data.adjacencies {
            if province.same_as(&adj.0) {
                println!("  {}", adj.1.name);
                count += 1;
            } else if province.same_as(&adj.1) {
                println!("  {}", adj.0.name);
                count += 1;
            }
        }
        println!("COUNT {}\n", count);
    }
}

enum Lobby {
    Open { users: Vec<String> },
    InGame,
}

struct LobbyUser {
    id: u32,
    name: String,
    role_preferences: Vec<String>,
}

mod game_data {
    pub struct GameData {
        pub nations: Vec<Nation>,
        pub provinces: Vec<Province>,
        pub adjacencies: Vec<(Province, Province)>,
    }

    impl GameData {
        pub fn default() -> Self {
            // Austrian provinces.
            let boh = Province {
                name: "Bohemia".to_owned(),
                abbreviation: "Boh".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let bud = Province {
                name: "Budapest".to_owned(),
                abbreviation: "Bud".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let gal = Province {
                name: "Galicia".to_owned(),
                abbreviation: "Gal".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let tri = Province {
                name: "Trieste".to_owned(),
                abbreviation: "Tri".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let tyr = Province {
                name: "Tyrolia".to_owned(),
                abbreviation: "Tyr".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let vie = Province {
                name: "Vienna".to_owned(),
                abbreviation: "Vie".to_owned(),
                terrain: ProvinceType::Land(None),
            };

            // English provinces.
            let cly = Province {
                name: "Clyde".to_owned(),
                abbreviation: "Cly".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let edi = Province {
                name: "Edinburgh".to_owned(),
                abbreviation: "Edi".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let lvp = Province {
                name: "Liverpool".to_owned(),
                abbreviation: "Lvp".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let lon = Province {
                name: "London".to_owned(),
                abbreviation: "Lon".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let wal = Province {
                name: "Wales".to_owned(),
                abbreviation: "Wal".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let yor = Province {
                name: "Yorkshire".to_owned(),
                abbreviation: "Yor".to_owned(),
                terrain: ProvinceType::Land(None),
            };

            // French provinces.
            let bre = Province {
                name: "Brest".to_owned(),
                abbreviation: "Bre".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let bur = Province {
                name: "Burgundy".to_owned(),
                abbreviation: "Bur".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let gas = Province {
                name: "Gascony".to_owned(),
                abbreviation: "Gas".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let mar = Province {
                name: "Marseilles".to_owned(),
                abbreviation: "Mar".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let par = Province {
                name: "Paris".to_owned(),
                abbreviation: "Par".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let pic = Province {
                name: "Picardy".to_owned(),
                abbreviation: "Pic".to_owned(),
                terrain: ProvinceType::Land(None),
            };

            // German provinces.
            let ber = Province {
                name: "Berlin".to_owned(),
                abbreviation: "Ber".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let kie = Province {
                name: "Kiel".to_owned(),
                abbreviation: "Kie".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let mun = Province {
                name: "Munich".to_owned(),
                abbreviation: "Mun".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let pru = Province {
                name: "Prussia".to_owned(),
                abbreviation: "Pru".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let ruh = Province {
                name: "Ruhr".to_owned(),
                abbreviation: "Ruh".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let sil = Province {
                name: "Silesia".to_owned(),
                abbreviation: "Sil".to_owned(),
                terrain: ProvinceType::Land(None),
            };

            // Italian provinces.
            let apu = Province {
                name: "Apulia".to_owned(),
                abbreviation: "Apu".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let nap = Province {
                name: "Naples".to_owned(),
                abbreviation: "Nap".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let pie = Province {
                name: "Piedmont".to_owned(),
                abbreviation: "Pie".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let rom = Province {
                name: "Rome".to_owned(),
                abbreviation: "Rom".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let tus = Province {
                name: "Tuscany".to_owned(),
                abbreviation: "Tus".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let ven = Province {
                name: "Venice".to_owned(),
                abbreviation: "Ven".to_owned(),
                terrain: ProvinceType::Land(None),
            };

            // Russian provinces.
            let fin = Province {
                name: "Finland".to_owned(),
                abbreviation: "Fin".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let lvn = Province {
                name: "Livonia".to_owned(),
                abbreviation: "Lvn".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let mos = Province {
                name: "Moscow".to_owned(),
                abbreviation: "Mos".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let sev = Province {
                name: "Sevastopol".to_owned(),
                abbreviation: "Sev".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let stp_nc = Province {
                name: "St. Petersburg".to_owned(),
                abbreviation: "Stp".to_owned(),
                terrain: ProvinceType::Land(Some(Coast::North)),
            };
            let stp_sc = Province {
                name: "St. Petersburg".to_owned(),
                abbreviation: "Stp".to_owned(),
                terrain: ProvinceType::Land(Some(Coast::South)),
            };
            let ukr = Province {
                name: "Ukraine".to_owned(),
                abbreviation: "Ukr".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let war = Province {
                name: "Warsaw".to_owned(),
                abbreviation: "War".to_owned(),
                terrain: ProvinceType::Land(None),
            };

            // Turkish provinces.
            let ank = Province {
                name: "Ankara".to_owned(),
                abbreviation: "Ank".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let arm = Province {
                name: "Armenia".to_owned(),
                abbreviation: "Arm".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let con = Province {
                name: "Constantinople".to_owned(),
                abbreviation: "Con".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let smy = Province {
                name: "Smyrna".to_owned(),
                abbreviation: "Smy".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let syr = Province {
                name: "Syria".to_owned(),
                abbreviation: "Syr".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let alb = Province {
                name: "Albania".to_owned(),
                abbreviation: "Alb".to_owned(),
                terrain: ProvinceType::Land(None),
            };

            // Neutral provinces.
            let bel = Province {
                name: "Belgium".to_owned(),
                abbreviation: "Bel".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let bul_ec = Province {
                name: "Bulgaria".to_owned(),
                abbreviation: "Bul".to_owned(),
                terrain: ProvinceType::Land(Some(Coast::East)),
            };
            let bul_sc = Province {
                name: "Bulgaria".to_owned(),
                abbreviation: "Bul".to_owned(),
                terrain: ProvinceType::Land(Some(Coast::South)),
            };
            let den = Province {
                name: "Denmark".to_owned(),
                abbreviation: "Den".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let gre = Province {
                name: "Greece".to_owned(),
                abbreviation: "Gre".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let hol = Province {
                name: "Holland".to_owned(),
                abbreviation: "Hol".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let nwy = Province {
                name: "Norway".to_owned(),
                abbreviation: "Nwy".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let naf = Province {
                name: "North Africa".to_owned(),
                abbreviation: "Naf".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let por = Province {
                name: "Portugal".to_owned(),
                abbreviation: "Por".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let rum = Province {
                name: "Rumania".to_owned(),
                abbreviation: "Rum".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let ser = Province {
                name: "Serbia".to_owned(),
                abbreviation: "Ser".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let spa_nc = Province {
                name: "Spain".to_owned(),
                abbreviation: "Spa".to_owned(),
                terrain: ProvinceType::Land(Some(Coast::North)),
            };
            let spa_sc = Province {
                name: "Spain".to_owned(),
                abbreviation: "Spa".to_owned(),
                terrain: ProvinceType::Land(Some(Coast::South)),
            };
            let swe = Province {
                name: "Sweden".to_owned(),
                abbreviation: "Swe".to_owned(),
                terrain: ProvinceType::Land(None),
            };
            let tun = Province {
                name: "Tunis".to_owned(),
                abbreviation: "Tun".to_owned(),
                terrain: ProvinceType::Land(None),
            };

            // Water provinces.
            let adr = Province {
                name: "Adriatic Sea".to_owned(),
                abbreviation: "ADR".to_owned(),
                terrain: ProvinceType::Water,
            };
            let aeg = Province {
                name: "Aegean Sea".to_owned(),
                abbreviation: "AEG".to_owned(),
                terrain: ProvinceType::Water,
            };
            let bal = Province {
                name: "Baltic Sea".to_owned(),
                abbreviation: "BAL".to_owned(),
                terrain: ProvinceType::Water,
            };
            let bar = Province {
                name: "Barents Sea".to_owned(),
                abbreviation: "BAR".to_owned(),
                terrain: ProvinceType::Water,
            };
            let bla = Province {
                name: "Black Sea".to_owned(),
                abbreviation: "BLA".to_owned(),
                terrain: ProvinceType::Water,
            };
            let eas = Province {
                name: "Eastern Mediterranean".to_owned(),
                abbreviation: "EAS".to_owned(),
                terrain: ProvinceType::Water,
            };
            let eng = Province {
                name: "English Channel".to_owned(),
                abbreviation: "ENG".to_owned(),
                terrain: ProvinceType::Water,
            };
            let bot = Province {
                name: "Gulf of Bothnia".to_owned(),
                abbreviation: "BOT".to_owned(),
                terrain: ProvinceType::Water,
            };
            let lyo = Province {
                name: "Gulf of Lyon".to_owned(),
                abbreviation: "LYO".to_owned(),
                terrain: ProvinceType::Water,
            };
            let hel = Province {
                name: "Helgoland Bight".to_owned(),
                abbreviation: "HEL".to_owned(),
                terrain: ProvinceType::Water,
            };
            let ion = Province {
                name: "Ionian Sea".to_owned(),
                abbreviation: "ION".to_owned(),
                terrain: ProvinceType::Water,
            };
            let iri = Province {
                name: "Irish Sea".to_owned(),
                abbreviation: "IRI".to_owned(),
                terrain: ProvinceType::Water,
            };
            let mao = Province {
                name: "Mid-Atlantic Ocean".to_owned(),
                abbreviation: "MAO".to_owned(),
                terrain: ProvinceType::Water,
            };
            let nth = Province {
                name: "North Sea".to_owned(),
                abbreviation: "NTH".to_owned(),
                terrain: ProvinceType::Water,
            };
            let nao = Province {
                name: "North Atlantic Ocean".to_owned(),
                abbreviation: "NAO".to_owned(),
                terrain: ProvinceType::Water,
            };
            let nwg = Province {
                name: "Norwegian Sea".to_owned(),
                abbreviation: "NWG".to_owned(),
                terrain: ProvinceType::Water,
            };
            let ska = Province {
                name: "Skagerrak".to_owned(),
                abbreviation: "SKA".to_owned(),
                terrain: ProvinceType::Water,
            };
            let tys = Province {
                name: "Tyrrhenian Sea".to_owned(),
                abbreviation: "TYS".to_owned(),
                terrain: ProvinceType::Water,
            };
            let wes = Province {
                name: "Western Mediterranean".to_owned(),
                abbreviation: "WES".to_owned(),
                terrain: ProvinceType::Water,
            };

            // Nations.
            let austria = Nation {
                name: "Austria-Hungary".to_owned(),
                home_supply_centers: vec![
                    (tri.clone(), Unit::Fleet),
                    (vie.clone(), Unit::Army),
                    (bud.clone(), Unit::Army),
                ],
            };
            let england = Nation {
                name: "England".to_owned(),
                home_supply_centers: vec![
                    (edi.clone(), Unit::Fleet),
                    (lvp.clone(), Unit::Army),
                    (lon.clone(), Unit::Fleet),
                ],
            };
            let france = Nation {
                name: "France".to_owned(),
                home_supply_centers: vec![
                    (bre.clone(), Unit::Fleet),
                    (par.clone(), Unit::Army),
                    (ruh.clone(), Unit::Army),
                ],
            };
            let germany = Nation {
                name: "Germany".to_owned(),
                home_supply_centers: vec![
                    (kie.clone(), Unit::Fleet),
                    (ber.clone(), Unit::Army),
                    (mun.clone(), Unit::Army),
                ],
            };
            let italy = Nation {
                name: "Italy".to_owned(),
                home_supply_centers: vec![
                    (nap.clone(), Unit::Fleet),
                    (rom.clone(), Unit::Army),
                    (ven.clone(), Unit::Army),
                ],
            };
            let russia = Nation {
                name: "Russia".to_owned(),
                home_supply_centers: vec![
                    (stp_sc.clone(), Unit::Fleet),
                    (mos.clone(), Unit::Army),
                    (war.clone(), Unit::Army),
                ],
            };
            let turkey = Nation {
                name: "Turkey".to_owned(),
                home_supply_centers: vec![
                    (ank.clone(), Unit::Fleet),
                    (con.clone(), Unit::Army),
                    (smy.clone(), Unit::Army),
                ],
            };

            let adjacencies = vec![
                (boh.clone(), mun.clone()),
                (boh.clone(), sil.clone()),
                (boh.clone(), gal.clone()),
                (boh.clone(), vie.clone()),
                (boh.clone(), tyr.clone()),
                (bud.clone(), vie.clone()),
                (bud.clone(), gal.clone()),
                (bud.clone(), rum.clone()),
                (bud.clone(), ser.clone()),
                (bud.clone(), tri.clone()),
                (gal.clone(), sil.clone()),
                (gal.clone(), war.clone()),
                (gal.clone(), ukr.clone()),
                (gal.clone(), rum.clone()),
                (gal.clone(), vie.clone()),
                (tri.clone(), ven.clone()),
                (tri.clone(), tyr.clone()),
                (tri.clone(), vie.clone()),
                (tri.clone(), ser.clone()),
                (tri.clone(), alb.clone()),
                (tri.clone(), adr.clone()),
                (tyr.clone(), mun.clone()),
                (tyr.clone(), vie.clone()),
                (tyr.clone(), ven.clone()),
                (tyr.clone(), pie.clone()),
                (cly.clone(), lvp.clone()),
                (cly.clone(), nao.clone()),
                (cly.clone(), nwg.clone()),
                (cly.clone(), edi.clone()),
                (edi.clone(), nwg.clone()),
                (edi.clone(), lvp.clone()),
                (edi.clone(), yor.clone()),
                (edi.clone(), nth.clone()),
                (lvp.clone(), nao.clone()),
                (lvp.clone(), yor.clone()),
                (lvp.clone(), wal.clone()),
                (lvp.clone(), iri.clone()),
                (lon.clone(), yor.clone()),
                (lon.clone(), nth.clone()),
                (lon.clone(), eng.clone()),
                (lon.clone(), wal.clone()),
                (wal.clone(), yor.clone()),
                (wal.clone(), eng.clone()),
                (wal.clone(), iri.clone()),
                (yor.clone(), nth.clone()),
                (bre.clone(), eng.clone()),
                (bre.clone(), pic.clone()),
                (bre.clone(), par.clone()),
                (bre.clone(), gas.clone()),
                (bre.clone(), mao.clone()),
                (bur.clone(), bel.clone()),
                (bur.clone(), ruh.clone()),
                (bur.clone(), mun.clone()),
                (bur.clone(), mar.clone()),
                (bur.clone(), gas.clone()),
                (bur.clone(), par.clone()),
                (bur.clone(), pic.clone()),
                (gas.clone(), par.clone()),
                (gas.clone(), mar.clone()),
                (gas.clone(), spa_nc.clone()),
                (gas.clone(), mao.clone()),
                (mar.clone(), pie.clone()),
                (mar.clone(), lyo.clone()),
                (mar.clone(), spa_sc.clone()),
                (par.clone(), pic.clone()),
                (pic.clone(), eng.clone()),
                (pic.clone(), bel.clone()),
                (ber.clone(), kie.clone()),
                (ber.clone(), bal.clone()),
                (ber.clone(), pru.clone()),
                (ber.clone(), sil.clone()),
                (ber.clone(), mun.clone()),
                (kie.clone(), hol.clone()),
                (kie.clone(), hel.clone()),
                (kie.clone(), den.clone()),
                (kie.clone(), bal.clone()),
                (kie.clone(), mun.clone()),
                (kie.clone(), ruh.clone()),
                (mun.clone(), ruh.clone()),
                (mun.clone(), sil.clone()),
                (pru.clone(), bal.clone()),
                (pru.clone(), lvn.clone()),
                (pru.clone(), war.clone()),
                (pru.clone(), sil.clone()),
                (ruh.clone(), bel.clone()),
                (ruh.clone(), hol.clone()),
                (sil.clone(), war.clone()),
                (apu.clone(), adr.clone()),
                (apu.clone(), ion.clone()),
                (apu.clone(), nap.clone()),
                (apu.clone(), rom.clone()),
                (apu.clone(), ven.clone()),
                (nap.clone(), ion.clone()),
                (nap.clone(), tys.clone()),
                (nap.clone(), rom.clone()),
                (pie.clone(), ven.clone()),
                (pie.clone(), tus.clone()),
                (pie.clone(), lyo.clone()),
                (rom.clone(), ven.clone()),
                (rom.clone(), tys.clone()),
                (tus.clone(), ven.clone()),
                (tus.clone(), rom.clone()),
                (tus.clone(), tys.clone()),
                (tus.clone(), lyo.clone()),
                (ven.clone(), adr.clone()),
                (fin.clone(), bot.clone()),
                (fin.clone(), swe.clone()),
                (fin.clone(), nwy.clone()),
                (fin.clone(), stp_sc.clone()),
                (lvn.clone(), bal.clone()),
                (lvn.clone(), bot.clone()),
                (lvn.clone(), stp_sc.clone()),
                (lvn.clone(), mos.clone()),
                (lvn.clone(), war.clone()),
                (mos.clone(), stp_sc.clone()),
                (mos.clone(), sev.clone()),
                (mos.clone(), ukr.clone()),
                (mos.clone(), war.clone()),
                (sev.clone(), arm.clone()),
                (sev.clone(), bla.clone()),
                (sev.clone(), rum.clone()),
                (sev.clone(), ukr.clone()),
                (stp_nc.clone(), bar.clone()),
                (stp_nc.clone(), nwy.clone()),
                (stp_sc.clone(), bot.clone()),
                (ukr.clone(), rum.clone()),
                (ukr.clone(), war.clone()),
                (ank.clone(), bla.clone()),
                (ank.clone(), arm.clone()),
                (ank.clone(), smy.clone()),
                (ank.clone(), con.clone()),
                (arm.clone(), bla.clone()),
                (arm.clone(), smy.clone()),
                (arm.clone(), syr.clone()),
                (con.clone(), bul_ec.clone()),
                (con.clone(), bul_sc.clone()),
                (con.clone(), bla.clone()),
                (con.clone(), smy.clone()),
                (con.clone(), aeg.clone()),
                (smy.clone(), aeg.clone()),
                (smy.clone(), eas.clone()),
                (smy.clone(), syr.clone()),
                (syr.clone(), eas.clone()),
                (alb.clone(), ser.clone()),
                (alb.clone(), adr.clone()),
                (alb.clone(), ion.clone()),
                (alb.clone(), gre.clone()),
                (bel.clone(), eng.clone()),
                (bel.clone(), nth.clone()),
                (bel.clone(), hol.clone()),
                (bul_ec.clone(), rum.clone()),
                (bul_ec.clone(), bla.clone()),
                (bul_sc.clone(), aeg.clone()),
                (bul_sc.clone(), gre.clone()),
                (bul_sc.clone(), ser.clone()),
                (den.clone(), nth.clone()),
                (den.clone(), ska.clone()),
                (den.clone(), bal.clone()),
                (den.clone(), hel.clone()),
                (den.clone(), swe.clone()),
                (gre.clone(), ser.clone()),
                (gre.clone(), ion.clone()),
                (gre.clone(), aeg.clone()),
                (hol.clone(), nth.clone()),
                (hol.clone(), hel.clone()),
                (nwy.clone(), swe.clone()),
                (nwy.clone(), nth.clone()),
                (nwy.clone(), ska.clone()),
                (nwy.clone(), nwg.clone()),
                (nwy.clone(), bar.clone()),
                (naf.clone(), tun.clone()),
                (naf.clone(), mao.clone()),
                (naf.clone(), wes.clone()),
                (por.clone(), mao.clone()),
                (por.clone(), spa_nc.clone()),
                (por.clone(), spa_sc.clone()),
                (rum.clone(), bla.clone()),
                (rum.clone(), ser.clone()),
                (spa_nc.clone(), mao.clone()),
                (spa_sc.clone(), mao.clone()),
                (spa_sc.clone(), wes.clone()),
                (spa_sc.clone(), lyo.clone()),
                (swe.clone(), ska.clone()),
                (swe.clone(), bal.clone()),
                (swe.clone(), bot.clone()),
                (tun.clone(), wes.clone()),
                (tun.clone(), tys.clone()),
                (tun.clone(), ion.clone()),
                (adr.clone(), ion.clone()),
                (aeg.clone(), eas.clone()),
                (aeg.clone(), ion.clone()),
                (bal.clone(), bot.clone()),
                (bar.clone(), nwg.clone()),
                (eas.clone(), ion.clone()),
                (eng.clone(), iri.clone()),
                (eng.clone(), mao.clone()),
                (eng.clone(), nth.clone()),
                (lyo.clone(), tys.clone()),
                (lyo.clone(), wes.clone()),
                (hel.clone(), nth.clone()),
                (ion.clone(), tys.clone()),
                (iri.clone(), nao.clone()),
                (iri.clone(), mao.clone()),
                (mao.clone(), nao.clone()),
                (mao.clone(), wes.clone()),
                (nao.clone(), nwg.clone()),
                (nth.clone(), nwg.clone()),
                (nth.clone(), ska.clone()),
                (tys.clone(), wes.clone()),
            ];

            let provinces = vec![
                boh, bud, gal, tri, tyr, vie, cly, edi, lvp, lon, wal, yor, bre, bur, gas, mar,
                par, pic, ber, kie, mun, pru, ruh, sil, apu, nap, pie, rom, tus, ven, fin, lvn,
                mos, sev, stp_sc, stp_nc, ukr, war, ank, arm, con, smy, syr, alb, bel, bul_sc,
                bul_ec, den, gre, hol, nwy, naf, por, rum, ser, spa_sc, spa_nc, swe, tun, adr, aeg,
                bal, bar, bla, eas, eng, bot, lyo, hel, ion, iri, mao, nao, nth, nwg, ska, tys,
                wes,
            ];
            let nations = vec![austria, england, france, germany, italy, russia, turkey];

            GameData {
                nations,
                provinces,
                adjacencies,
            }
        }
    }

    #[derive(Clone)]
    pub struct Nation {
        pub name: String,
        pub home_supply_centers: Vec<(Province, Unit)>,
    }

    #[derive(Clone)]
    pub struct Province {
        pub name: String,
        pub abbreviation: String,
        pub terrain: ProvinceType,
    }

    impl Province {
        pub fn same_as(&self, other: &Province) -> bool {
            self.abbreviation == other.abbreviation
        }

        pub fn same_coast_as(&self, other: &Province) -> bool {
            self.abbreviation == other.abbreviation && self.terrain == other.terrain
        }
    }

    #[derive(Clone, Copy, PartialEq, Eq)]
    pub enum ProvinceType {
        Land(Option<Coast>),
        Water,
    }

    #[derive(Clone, Copy, PartialEq, Eq)]
    pub enum Coast {
        North,
        South,
        East,
        West,
    }

    #[derive(Clone, Copy, PartialEq, Eq)]
    pub enum Unit {
        Army,
        Fleet,
    }
}
