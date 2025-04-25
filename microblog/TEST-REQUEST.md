# Comandos grpcurl para Testar os Microserviços

## UserService (localhost:8081)

### Register (Criar usuário)
```bash
grpcurl -plaintext -d '{
  "username": "valentina",
  "email": "valentina@example.com",
  "password": "123456"
}' localhost:8081 user.UserService/Register
```

### Login (Autenticar usuário)
```bash
grpcurl -plaintext -d '{
  "email": "valentina@example.com",
  "password": "123456"
}' localhost:8081 user.UserService/Login
```

### GetUser (Buscar usuário por ID)
```bash
grpcurl -plaintext -d '{
  "id": "<user_id>"
}' localhost:8081 user.UserService/GetUser
```

### ValidateToken (Validar JWT)
```bash
grpcurl -plaintext -d '{
  "token": "<jwt_token>"
}' localhost:8081 user.UserService/ValidateToken
```

---

## PostService (localhost:8082)

### CreatePost (Criar post)
```bash
grpcurl -plaintext -d '{
  "user_id": "user456",
  "content": "Post final da streak!"
}' localhost:8082 post.PostService/CreatePost
```

### ListPosts (Listar posts de um usuário)
```bash
grpcurl -plaintext -d '{
  "user_id": "user456",
  "limit": 10
}' localhost:8082 post.PostService/ListPosts
```

---

## FollowService (localhost:8083)

### Follow (Seguir usuário)
```bash
grpcurl -plaintext -d '{
  "follower_id": "user123",
  "followee_id": "user456"
}' localhost:8083 follow.FollowService/Follow
```

### ListFollowers (Listar followers de um usuário)
```bash
grpcurl -plaintext -d '{
  "user_id": "user456"
}' localhost:8083 follow.FollowService/ListFollowers
```

---

## TimelineService (localhost:8084)

### GetTimeline (Buscar timeline de um usuário)
```bash
grpcurl -plaintext -d '{
  "user_id": "user123",
  "limit": 10
}' localhost:8084 timeline.TimelineService/GetTimeline
```

