# Доска объявлений


## Frontend

### Установка

```sh
cd client
npm install
```

### Запуск

```sh
npm run dev
```

## Backend 
### Запуск
```sh
cd server
go run . localhost:8080
```


## api

`/bulletin_board` - работа с доской объявлений

```ts
type ServerResponse = {
  AnswerType: string;
  Value: string[];
};
```

* AnswerType - Либо LIST, 
