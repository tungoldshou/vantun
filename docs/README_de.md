# VANTUN - Sicheres Tunnelprotokoll der nächsten Generation

VANTUN ist ein cutting-edge, Hochleistungs-Tunnelprotokoll basierend auf QUIC, das entwickelt wurde, um außergewöhnliche Netzwerkleistung, Sicherheit und Zuverlässigkeit zu bieten. Als Next-Generation-Lösung definiert VANTUN mit seiner innovativen Architektur und fortschrittlichen Funktionen neu, was im Netzwerktunneling möglich ist.

## Hauptvorteile

### 🔒 Sicherheit auf Unternehmensebene
- **Sicheres Handshake und Sitzungsverhandlung**: Durchgeführt über dedizierten Kontroll-Stream für Verbindungssicherheit

### ⚡ Außergewöhnliche Leistung
- **Mehrere Logische Stream-Typen**: Optimierte interaktive, Bulk- und Telemetrie-Streams für verschiedene Geschäftsszenarien
- **Multipath**: Intelligente Nutzung mehrerer Netzwerkpfade für dramatisch verbesserte Geschwindigkeit und Verbindungsstabilität

### 🛡️ Unübertroffene Zuverlässigkeit
- **Forward Error Correction (FEC)**: Fortschrittliche Fehlerkorrektur gewährleistet Datenintegrität selbst bei instabilen Netzwerkbedingungen
- **Hybride Congestion Control**: Innovativer Hybrid-Algorithmus, der QUIC CC mit Token-Bucket-Rate-Limiting kombiniert für optimale Ressourcennutzung

### 🌐 Privatsphärenschutz
- **Steckbarer Obfuskations-Modul**: Fortschrittliche Traffic-Obfuskation lässt Traffic wie normalen HTTP/3 erscheinen, um Netzwerkprüfung effektiv zu umgehen

### 🚀 Einfache Bereitstellung
- **Minimale Client/Server**: Befehlszeilen-`client`- und `server`-Programme für schnelle Bereitstellung und einfachen Gebrauch

## Technologie-Architektur

VANTUN nutzt branchenführende Technologien, um außergewöhnliche Leistung und Zuverlässigkeit zu bieten:

- **Sprache**: Go - Hochleistungsfähige, moderne Programmiersprache mit Nebenläufigkeit
- **Kern-Bibliothek**: `quic-go` - Branchenführende QUIC-Protokoll-Implementierung
- **Serialisierung**: `github.com/fxamacker/cbor` - Effiziente CBOR-Kodierung, kompakter als JSON
- **FEC**: `github.com/klauspost/reedsolomon` - Hochleistungsfähiger Reed-Solomon-Kodierungsalgorithmus
- **CLI**: `cobra/viper` - Mächtige Befehlszeilenschnittstelle und Konfigurationsmanagement

## Schnellstart

Bringen Sie VANTUN in wenigen Minuten zum Laufen:

1. **Repository Klonen**: `git clone <repository-url>`
2. **Erstellen**: `go build -o bin/vantun cmd/main.go`
3. **Konfigurieren**: Erstellen Sie die `config.json`-Konfigurationsdatei
4. **Ausführen**: Server und Client starten

Für detaillierte Schritte und Konfigurationsanweisungen lesen Sie bitte die [Demo-Anleitung](DEMOGUIDE_de.md).

## Projektstruktur

```
vantun/
├── cmd/              # Befehlszeilenprogramm-Einstieg
├── internal/
│   ├── cli/          # CLI-Konfigurationsmanagement
│   └── core/         # Kern-Protokoll-Implementierung
├── docs/             # Dokumentation
├── go.mod            # Go-Modul-Definition
└── README.md         # Projektdokumentation
```

## Architektur-Highlights

### 🔧 Intelligente Protokoll-Engine
Die Kern-Protokoll-Engine implementiert effiziente Sitzungsverhandlung und Kontroll-Stream-Management für sichere und stabile Verbindungen.

### 📊 Adaptive FEC-Technologie
Forward Error Correction basierend auf Reed-Solomon-Kodierung, die dynamisch Korrekturstrategien basierend auf Netzwerkbedingungen anpasst.

### 🔄 Intelligente Multipath-Übertragung
Innovatives Pfad-Management und Load-Balancing, das alle verfügbaren Netzwerkpfade vollständig für Redundanz und verbesserten Durchsatz nutzt.

### 📈 Hybride Congestion Control
Hybrid-Algorithmus, der zugrunde liegende QUIC-Congestion-Control mit oberer Schicht Token-Bucket kombiniert für optimale Ressourcennutzung.

### 🎭 Fortgeschrittene Traffic-Obfuskation
HTTP/3-ähnliche Traffic-Obfuskation und intelligente Daten-Auffüllung, um Netzwerkprüfung effektiv zu umgehen und Benutzer-Privatsphäre zu schützen.

### 📊 Echtzeit-Telemetrie-System
Umfassende Performance-Datensammlung und Echtzeit-Überwachung für Netzwerk-Optimierung und Fehlerbehebung.

## Qualitätssicherung

VANTUN übernimmt strenge Teststandards, um Code-Qualität und System-Stabilität sicherzustellen:

- **Umfassende Unit-Tests**: Abdeckung aller Kern-Funktionsmodule
- **Integrationstests**: Validierung der Komponenten-Zusammenarbeit
- **Performance-Tests**: Sicherstellung außergewöhnlicher Performance unter verschiedenen Netzwerkbedingungen
- **Stresstests**: Validierung der Stabilität unter hoher Last

Alle Tests ausführen:

```bash
go test -v ./internal/core/...
```

## Lizenz

VANTUN ist unter der MIT-Lizenz lizenziert, einer permissiven Open-Source-Lizenz, die freie Nutzung, Kopierung, Modifikation und Verteilung der Software erlaubt, während Urheberrechts- und Lizenznotizen beibehalten werden.

---

*© 2025 VANTUN-Projekt. Alle Rechte vorbehalten.*