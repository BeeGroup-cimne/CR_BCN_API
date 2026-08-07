# 🚀 CRAPI – Climate Ready Barcelona API

🌍 *Smart data for a resilient future*

---

## 📘 Description

**CRAPI (Climate Ready Barcelona API)** is a RESTful API providing access to climate, urban, and resilience data for the city of Barcelona.
Key features include:
- Access to climate and urban sustainability indicators
- Data segmented by territorial zones (districts, neighborhoods, RTAs, etc.)
- Create dynamic vulnerability indices based on user-defined groups of indicators
- Evaluate the available climate shelters for each building based on schedules, prices, features, and equipment.
- Prepared data for climate resilience analysis

---

## 🛠️ Technologies

- Go (Golang) `v1.23+`
- Swagger (OpenAPI) for interactive documentation
- MongoDB / Neo4J / HBase
- Docker 
- JWT (for authentication)
- K8s 

---

## 🚦 Main Endpoints

Explore the interactive API documentation via Swagger UI:

🔗 [`/swagger/index.html`](https://crbcn-go-api.kubernetes.tech.beegroup-cimne.com/swagger/index.html#/)

---

## ▶️ How to Run

To start the API server, run the main Go application located at 🔗 [`cmd/main.go`](cmd/main.go)
Make sure you have Go installed and configured in your environment. The server will start on the configured port `:8090`.


## ✨ How to Contribute

Got ideas or improvements? Contributions are welcome!

1. Fork the repository
2. Create a new branch: `git checkout -b feature/new-feature`
3. Commit your changes
4. Open a Pull Request explaining your contribution

---

## 👨‍💻 Authors

This project was developed by:

- **Jose Manuel Broto** – jmbroto@cimne.upc.edu
- **Gerard Mor** – gmor@cimne.upc.edu
- **Andreu Marí** – amari@cimne.upc.edu

---

## 📫 Contact

📧 Questions or suggestions?  
Email us!

---
## 🧾 License

This project is licensed under the **EUPL License**. See the [LICENSE](./LICENSE) file for details.

---

#docker buildx build --platform linux/amd64,linux/arm64 --push -t docker.tech.beegroup-cimne.com/crbcn_api/go-api . --provenance=false
docker buildx build --platform linux/amd64,linux/arm64 --push -t 853583158095.dkr.ecr.eu-west-1.amazonaws.com/ecr-crbcn-api . --provenance=false
kubectl create secret generic crbcn-api-secret --from-file=.env=.env -n crbcn --dry-run=client -o yaml | kubectl apply -f -