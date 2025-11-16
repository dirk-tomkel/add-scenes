# Kdenlive Scene Importer (Go)

# Kdenlive Scene Importer (`add-scenes`)

Ein einfaches Go-Tool, das **automatisch Szenen aus einer `pyscenedetect`-CSV-Datei** in eine **Kdenlive-Projektdatei** (`.kdenlive`) importiert – ideal für schnelle Schnittvorbereitung.

> **GitHub:** [https://github.com/dirk-tomkel/add-scenes](https://github.com/dirk-tomkel/add-scenes)

---

## Was macht das Tool?

1. Öffnet deine `.kdenlive`-Projektdatei.
2. Sucht alle **Playlisten außer `main_bin`**.
3. **Löscht** alle vorhandenen Einträge in diesen Playlisten.
4. Liest die Szenen aus der CSV-Datei (ab Zeile 3).
5. Fügt **neue `<entry>`-Einträge** mit `in`, `out` und dem passenden `producer` hinzu.
6. Speichert die neue Datei als:  
   `scenes-<originalname>.kdenlive`

---

## Voraussetzungen

- Eine `.kdenlive`-Projektdatei
- Eine CSV-Datei von `pyscenedetect` (mit `--dump-csv`)

> **Kein Go-Wissen nötig** – siehe Installation für Nicht-Entwickler!

---

## Installation

```bash
go install github.com/dirk-tomkel/add-scenes@latest