```mermaid
graph LR;
    A[user] -- GET --> B[localhost:8080 <br>go server]
    B -- GET /, /assets--> C[localhost:3000 <br> node server]
    C --> D["dist"]
    C -- react render to string --> E["App"]
    F[node.ts] -- tsc --> C
    G[index.html template] -- vite--> D

```