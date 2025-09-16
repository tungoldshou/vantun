# VANTUN - Protocole de Tunnel Sécurisé de Nouvelle Génération

VANTUN est un protocole de tunnel de pointe à haute performance basé sur QUIC, conçu pour offrir des performances réseau exceptionnelles, une sécurité et une fiabilité optimales. En tant que solution de nouvelle génération, VANTUN redéfinit ce qui est possible dans le tunnel réseau avec son architecture innovante et ses fonctionnalités avancées.

## Avantages Clés

### 🔒 Sécurité de Niveau Entreprise
- **Poignée de Main Sécurisée et Négociation de Session**: Effectuées via un flux de contrôle dédié pour la sécurité de la connexion

### ⚡ Performance Exceptionnelle
- **Types de Flux Logiques Multiples**: Flux interactifs, volumineux et de télémétrie optimisés pour différents scénarios d'affaires
- **Multipath**: Utilisation intelligente de plusieurs chemins réseau pour une vitesse et une stabilité de connexion drastiquement améliorées

### 🛡️ Fiabilité Inégalée
- **Correction d'Erreur Anticipée (FEC)**: Technologie de correction d'erreur avancée garantissant l'intégrité des données même dans des conditions réseau instables
- **Contrôle de Congestion Hybride**: Algorithme innovant combinant le contrôle de congestion QUIC avec la limitation de débit par seau à jetons pour une utilisation optimale des ressources

### 🌐 Protection de la Vie Privée
- **Module d'Obfuscation Plugable**: Obfuscation de trafic avancée faisant apparaître le trafic comme du HTTP/3 normal, évitant efficacement l'examen du réseau

### 🚀 Déploiement Facile
- **Client/Serveur Minimaux**: Programmes en ligne de commande `client` et `server` pour un déploiement rapide et une facilité d'utilisation

## Architecture Technologique

VANTUN exploite des technologies de pointe de l'industrie pour offrir ses performances et fiabilité exceptionnelles:

- **Langage**: Go - Langage de programmation moderne haute performance et concurrent
- **Bibliothèque Cœur**: `quic-go` - Implémentation de protocole QUIC leader dans l'industrie
- **Sérialisation**: `github.com/fxamacker/cbor` - Encodage CBOR efficace, plus compact que JSON
- **FEC**: `github.com/klauspost/reedsolomon` - Algorithme d'encodage Reed-Solomon haute performance
- **CLI**: `cobra/viper` - Interface en ligne de commande puissante et gestion de configuration

## Démarrage Rapide

Mettez VANTUN en route en quelques minutes:

1. **Cloner le Dépôt**: `git clone <repository-url>`
2. **Construire**: `go build -o bin/vantun cmd/main.go`
3. **Configurer**: Créer le fichier de configuration `config.json`
4. **Exécuter**: Démarrer le serveur et le client

Pour des instructions détaillées et la configuration, veuillez consulter le [Guide de Démo](DEMOGUIDE_fr.md).

## Structure du Projet

```
vantun/
├── cmd/              # Point d'entrée du programme en ligne de commande
├── internal/
│   ├── cli/          # Gestion de configuration CLI
│   └── core/         # Implémentation du protocole cœur
├── docs/             # Documentation
├── go.mod            # Définition du module Go
└── README.md         # Documentation du projet
```

## Points Forts de l'Architecture

### 🔧 Moteur de Protocole Intelligent
Le moteur de protocole cœur implémente une négociation de session efficace et une gestion de flux de contrôle pour des connexions sécurisées et stables.

### 📊 Technologie FEC Adaptative
Correction d'erreur anticipée basée sur l'encodage Reed-Solomon qui ajuste dynamiquement les stratégies de correction basées sur les conditions réseau.

### 🔄 Transmission Multipath Intelligente
Gestion de chemin innovante et équilibrage de charge qui utilise pleinement tous les chemins réseau disponibles pour la redondance et le débit amélioré.

### 📈 Contrôle de Congestion Hybride
Algorithme hybride combinant le contrôle de congestion QUIC sous-jacent avec le seau à jetons de couche supérieure pour une utilisation optimale des ressources.

### 🎭 Obfuscation de Trafic Avancée
Obfuscation de trafic de style HTTP/3 et remplissage de données intelligent pour éviter efficacement l'examen du réseau et protéger la vie privée de l'utilisateur.

### 📊 Système de Télémétrie en Temps Réel
Collecte de données de performance complète et surveillance en temps réel pour l'optimisation réseau et le dépannage.

## Assurance Qualité

VANTUN adopte des normes de test strictes pour garantir la qualité du code et la stabilité du système:

- **Tests Unitaires Complets**: Couvrant tous les modules fonctionnels cœur
- **Tests d'Intégration**: Validation de la collaboration des composants
- **Tests de Performance**: Garantie de performance exceptionnelle sous diverses conditions réseau
- **Tests de Stress**: Validation de la stabilité sous charge élevée

Exécuter tous les tests:

```bash
go test -v ./internal/core/...
```

## Licence

VANTUN est sous licence MIT, une licence open-source permissive qui permet l'utilisation, la copie, la modification et la distribution gratuites du logiciel tout en conservant les avis de droit d'auteur et de licence.

---

*© 2025 Projet VANTUN. Tous droits réservés.*