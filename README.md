# 🤖 Package aiyou

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.16-blue.svg)](https://golang.org/dl/)

Un SDK Go minimaliste et efficace pour interagir avec l'API AI.You. Conçu pour être simple d'utilisation tout en offrant une flexibilité maximale.

## ✨ Caractéristiques

- 🎯 Interface simple et intuitive
- 🤖 Deux fonctions principales :
  - Lister les modèles disponibles
  - Envoyer un message à un modèle ou un assistant, contrôle de la température
- 🔄 Support du streaming pour les réponses en temps réel
- ⚡ Gestion automatique des retries
- 🛠️ Options de configuration flexibles
- 🔍 Mode debug pour le développement
- 🌡️ Contrôle de la température des réponses

## 📦 Installation

```bash
go get github.com/n1neT10ne/aiyou
```

## 🚀 Utilisation

```go
package main

import (
    "fmt"
    "github.com/n1neT10ne/aiyou"
)

func main() {
    // Liste des modèles disponibles
    models, err := aiyou.ListModels(
        "your-token",
        aiyou.WithDebug(true),
    )
    if err != nil {
        panic(err)
    }
    for _, model := range models {
        fmt.Printf("- %s\n", model.Name)
    }

    // Configuration de base
    response, err := aiyou.Completion(
        "model-name",
        "your-token",
        "votre message",
    )
    if err != nil {
        panic(err)
    }
    fmt.Println(response)

    // Avec options
    response, err = aiyou.Completion(
        "model-name",
        "your-token",
        "votre message",
        aiyou.WithDebug(true),
        aiyou.WithTemperature(0.7),
        aiyou.WithSystemPrompt("prompt système"),
        aiyou.WithStream(true),
    )
}
```

## ⚙️ Options

Le package supporte plusieurs options de configuration :

```go
// Active le mode debug
WithDebug(debug bool)

// Définit la température pour la génération (0.0-2.0)
WithTemperature(temp float64)

// Configure les tentatives en cas d'erreur
WithRetry(maxRetries int, delay time.Duration)

// Définit le prompt système
WithSystemPrompt(prompt string)

// Active le mode streaming
WithStream(stream bool)
```

## 🔄 Mode Streaming

Le mode streaming permet de recevoir la réponse au fur et à mesure qu'elle est générée. Il est particulièrement utile pour les réponses longues ou pour afficher la réponse progressivement.

```go
response, err := aiyou.Completion(
    "model-name",
    "your-token",
    "votre message",
    aiyou.WithStream(true),
)
```

En mode streaming :
- La réponse est construite progressivement à partir des chunks reçus
- Chaque chunk contient une partie de la réponse finale
- Le mode debug affiche les chunks reçus et leur contenu

## ⚠️ Gestion des erreurs

Le package définit plusieurs types d'erreurs :

```go
var (
    ErrEmptyToken      = errors.New("token is required")
    ErrEmptyMessage    = errors.New("message is required")
    ErrInvalidToken    = errors.New("invalid token")
    ErrRateLimit       = errors.New("rate limit exceeded")
    ErrStreamCorrupted = errors.New("stream response corrupted")
)
```

## 🧪 Tests Unitaires

Le package inclut une suite complète de tests unitaires. Pour les exécuter, vous devez définir votre token AI.You dans la variable d'environnement `AIYOU_TEST_TOKEN` :

```bash
export AIYOU_TEST_TOKEN="votre-token"
go test ./...
```

Si la variable d'environnement n'est pas définie, les tests nécessitant une authentification seront automatiquement ignorés avec un message explicatif. Cela permet de :
- Éviter de stocker des tokens directement dans le code
- Faciliter l'intégration continue sans exposer de données sensibles
- Permettre à chaque développeur d'utiliser son propre token de test

## 🔗 CLI

Un outil en ligne de commande est disponible dans un projet séparé : [aiyou-cli](https://github.com/n1neT10ne/aiyou-cli). Cette interface en ligne de commande offre un moyen rapide et simple d'interagir avec l'API AI.You directement depuis votre terminal.

## 📄 Licence

Ce projet est sous licence MIT - voir le fichier [LICENSE](LICENSE) pour plus de détails.

## 👤 Auteur

**Cyrille BARTHELEMY**

* Github: [@n1neT10ne](https://github.com/n1neT10ne)
