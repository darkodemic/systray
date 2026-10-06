<!--
Copyright 2026 Darko Demić.
Licensed under the Apache License, Version 2.0; see LICENSE and NOTICE.
-->

# 0001 — Putanja Go modula: github.com/darkodemic/systray

- **Status:** Prihvaćeno 2026-10-07. Modul se zove `github.com/darkodemic/systray` umesto `fyne.io/systray`.
- **Datum:** 2026-10-07
- **Zamenjuje:** — / **Zamenjena sa:** —
- **Vlasnik:** Darko
- **Povezano:** `NOTICE` (poreklo i autorska prava forka), `AGENTS.md` (grane za fyne-io PR-ove kreću od upstream `master`-a)

## Kontekst

- Fork je od fyne-io/systray nasledio `module fyne.io/systray`. Domen `fyne.io` pripada Fyne projektu, pa fork pod tim imenom izgleda kao njihov paket.
- Sa tom putanjom `go get github.com/darkodemic/systray` ne radi: Go odbija modul jer on za sebe kaže da je `fyne.io/systray`. Fork se zato može koristiti samo preko `replace` direktive. gpwebcam ga koristi upravo tako: `replace fyne.io/systray => github.com/darkodemic/systray v1.12.3-0.20261006205618-9c45f672f861` (provereno 2026-10-07).
- README je opisivao integraciju sa Fyne toolkit-om i `fyne package`, a to za ovaj fork nije relevantno.

## Odluka

1. `go.mod` deklariše `module github.com/darkodemic/systray`, a svi interni import-i, primer i README koriste tu putanju.
2. Iz README-a su uklonjeni delovi specifični za Fyne: sekcija o Fyne aplikaciji, preporuke za `fyne package` i linkovi na developer.fyne.io i pkg.go.dev/fyne.io.
3. Poreklo ostaje zapisano. README (uvod i Credits), `NOTICE` i `AGENTS.md` i dalje pominju fyne-io/systray, jer je to atribucija, a ne brending.

## Posledice

Pozitivne:

- Fork se koristi direktno, `go get github.com/darkodemic/systray`, bez `replace` direktive.
- pkg.go.dev prikazuje dokumentaciju forka pod njegovim sopstvenim imenom.

Negativne:

- Ovo je izmena koja lomi korisnike. Ko importuje `fyne.io/systray` sa `replace` na fork mora, pri prelasku na verziju posle ove izmene, da promeni import-e na novu putanju i da ukloni `replace`. Starija verzija na koju je projekat već pinovan (gpwebcam na 9c45f67) radi i dalje.
- Pri preuzimanju izmena iz upstream-a nastaju konflikti u import linijama pet fajlova (`systray_unix.go`, `systray_menu_unix.go`, `systray_notifier_unix.go`, `systray_unix_test.go`, `example/main.go`) kad god upstream menja te linije.

Rizici:

- PR ka fyne-io napravljen iz grane zasnovane na `master`-u forka nosio bi novu putanju. Pravilo iz `AGENTS.md`, po kome takve grane kreću od upstream `master`-a, to već sprečava.

## Razmotrene alternative

- **Zadržati `fyne.io/systray`.** Imalo bi manje konflikata sa upstream-om, ali bi fork ostao upotrebljiv samo preko `replace` direktive i nosio bi tuđi domen.
- **Vanity putanja (npr. `darkodemic.com/systray`).** Bila bi nezavisna od GitHub-a, ali zahteva hostovanje `go-import` meta tagova, a za tim sada nema potrebe.

## Van opsega

- Prelazak gpwebcam-a na novu putanju.
- `Makefile` cilj `tag-changelog`, koji i dalje generiše changelog iz getlantern/systray.
