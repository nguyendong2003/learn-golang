# Restful API

## Gin Web Framework
<https://gin-gonic.com/en/docs/quickstart/>

## GORM
<https://gorm.io/docs/>

```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
```

## godotenv  (Read env file)
<https://github.com/joho/godotenv>

```bash
go get github.com/joho/godotenv
```


## Postman API
- Create item: http://localhost:8080/v1/items
Body:
{
    "title": "This is a new item 99",
    "description": "Item description 99",
    "status": "Done"
}

- Get item by id: http://localhost:8080/v1/items/:id

- Update item by id: http://localhost:8080/v1/items/:id
Body: 
{
    "title": "abc22",
    "description": "description 22",
    "status":"Done"
}

- Delete item by id: http://localhost:8080/v1/items/:id

- Get items: http://localhost:8080/v1/items?page=2&limit=5&status=Done


