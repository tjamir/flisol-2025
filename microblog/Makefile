PROTO_DIRS := user-service post-service follow-service timeline-service

generate-proto: proto-tools check-protoc
	@echo "Gerando arquivos .pb.go para todos os serviços com suporte a Go modules..."
	@for dir in $(PROTO_DIRS); do \
		echo "Gerando protos para $$dir..."; \
		protoc \
			--proto_path=$$dir/proto \
			--go_out=$$dir --go_opt=module=github.com/tjamir/flisol-2025/microblog/$$dir \
			--go-grpc_out=$$dir --go-grpc_opt=module=github.com/tjamir/flisol-2025/microblog/$$dir \
			$$dir/proto/*.proto; \
	done
	@echo "Geração finalizada!"

proto-tools:
	@echo "Verificando ferramentas protoc-gen-go e protoc-gen-go-grpc..."
	@command -v protoc-gen-go >/dev/null 2>&1 || go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@command -v protoc-gen-go-grpc >/dev/null 2>&1 || go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Ferramentas verificadas."

check-protoc:
	@command -v protoc >/dev/null 2>&1 || (echo "Erro: 'protoc' não encontrado no PATH. Por favor, instale o protocolo buffer compiler (protoc) antes de continuar." && exit 1)

clean-proto:
	@echo "Removendo arquivos gerados (*.pb.go) de todos os serviços..."
	@for dir in $(PROTO_DIRS); do \
		find $$dir/proto -type f \( -name "*.pb.go" -o -name "*_grpc.pb.go" \) -exec rm -v {} \;; \
	done
	@echo "Limpeza finalizada!"
