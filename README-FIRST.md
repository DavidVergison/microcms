# Fait :

Instantiation de les tables :
- Sites
- Artices

Création des primary adapters :
- PUT articles
- Stream des articles modifiés (poc pour un futur déclenchement de génération)

# Next step :
- doc 4C + structure de données
- grosse relecture du premier endpoint
- grosse relecture sur les tests
- verification de la complétude des tests du usecase (il en manque forcement)
- Sécurisation OAuth via Cognito (infra as code dans template.yaml)

# Prerequis :
- Un compte AWS configuré (pour l'instant je déploie sur mon compte perso pour tester)

Déployer :
```sh
make # build
sam deploy --guided # first deployment
sam deploy # next deployments
```

# Structure
- articles : usecases liés aux contenus
- dynamodb : secondary adapters concernant dynamodb
- https : primary adapters de type API
- stream : primary adapters de type stream
- tests : tooling sur les tests