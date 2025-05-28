# IPX - IP eXpression Language

**IPX (IP eXpression Language)** is a lightweight, expressive protocol for parsing and interpreting user input into structured data formats like JSON. Built with performance and simplicity in mind, IPX helps transform ambiguous or hard-to-read strings into clear logical statements.

<p align="center">
  <img src="https://img.shields.io/github/license/cakebaker1/ipx?style=for-the-badge" />
  <img src="https://img.shields.io/github/stars/cakebaker1/ipx?style=for-the-badge" />
  <img src="https://img.shields.io/github/forks/cakebaker1/ipx?style=for-the-badge" />
  <img src="https://img.shields.io/github/issues/cakebaker1/ipx?style=for-the-badge" />
</p>

---

## 🌐 About

Many modern services struggle with correct user input parsing, often breaking logical operations in specific scenarios. IPX was created to address these shortcomings, introducing a standardized and extensible syntax for communication between the user and the machine.

Originally developed for [**KriX**](https://krix.world), IPX is now a standalone protocol that anyone can use, modify, and extend across any programming language or platform.

> 💡 We encourage developers to experiment, adapt, and expand the language without the burden of strict copyright enforcement.

---

## 🎯 Goal

Convert user-friendly query strings into structured formats that are easy to parse and reason about programmatically.

**Example:**

```
something:"some_value"
```

Will be transformed into:

```json
{
  "type": "MATCH",
  "key": "something",
  "op": ":",
  "value": "some_value"
}
```

This JSON structure can then be passed through any processing pipeline, rule engine, or search filter.

---

## ⚙️ Supported Operators


| Operator | Description            |
|----------|------------------------|
| :        | LIKE match (contains)  |
| :=       | Exact match (equals)   |
| &&       | Logical AND            |
| \|\|       | Logical OR             |
| !        | Logical NOT            |
| ()       | Grouping expressions   |
| []       | Array-like structures  |


---

## 🚀 Features

- 🔥 High performance
- 🛡️ Secure and safe parsing
- 🌍 Cross-platform compatibility
- 🧠 Simple and human-friendly syntax
- 🧩 Easy to extend and integrate
- ✅ Supports complex logical expressions
- 🛠️ Open to community improvements

---

## 📦 Install & Usage

_Coming soon: Installation instructions for different platforms and languages._

---

## 💡 Use Cases

- User input parsing
- Query building and filtering
- Search engine interpreters
- Rule-based engines
- Config string processors

---

## 📄 License

This project is licensed under the terms of the **MIT License**.  
See the [LICENSE](https://github.com/cakebaker1/ipx/blob/main/LICENSE) file for full details.

---

## 🤝 Contributing

We welcome contributions from the community!  
Feel free to submit issues, feature requests, or pull requests.

---

## 🔗 Links

- GitHub Repo: [https://github.com/cakebaker1/ipx](https://github.com/cakebaker1/ipx)
- Related project (KriX): [https://not_google.wow](https://google.com)

---

Made with ❤️ by the Krix_D3vs
