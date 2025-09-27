#!/bin/bash

# --- Configuração ---
# Diretório que será lido recursivamente
SOURCE_DIR="api"
# Arquivo de saída onde tudo será escrito. Ele será criado dentro do próprio diretório 'api'.
OUTPUT_FILE="api/estrutura_e_conteudo.txt"
# Lista de arquivos a serem ignorados
declare -a IGNORE_FILES=("go.mod" "go.sum" "smartlock")
# --- Fim da Configuração ---

# Validação: Verifica se o diretório de origem existe
if [ ! -d "$SOURCE_DIR" ]; then
    echo "Erro: O diretório de origem '$SOURCE_DIR' não foi encontrado."
    echo "Execute este script a partir da raiz do projeto ('Tranca-inteligente')."
    exit 1
fi

# Constrói os argumentos de exclusão para o comando find
FIND_ARGS=""
for file in "${IGNORE_FILES[@]}"; do
  FIND_ARGS+=" -not -name $file"
done

# Limpa/Cria o arquivo de saída para garantir que começamos com um arquivo vazio
echo "Criando/Limpando o arquivo de saída: $OUTPUT_FILE"
> "$OUTPUT_FILE"

echo "Iniciando a busca de arquivos... Ignorando: ${IGNORE_FILES[*]}"

# Encontra todos os arquivos (e não diretórios) dentro de SOURCE_DIR, aplicando as exclusões.
# O uso de `eval` aqui é seguro pois estamos controlando os argumentos que são construídos.
eval find "$SOURCE_DIR" -type f $FIND_ARGS -print0 | while IFS= read -r -d '' filepath; do
    echo "Processando: $filepath"

    # 1. Escreve o caminho relativo do arquivo no arquivo de saída
    echo "$filepath" >> "$OUTPUT_FILE"

    # 2. Escreve o conteúdo completo do arquivo logo abaixo
    cat "$filepath" >> "$OUTPUT_FILE"

    # 3. (Opcional, mas recomendado) Adiciona duas linhas em branco para separar
    #    visualmente o conteúdo dos diferentes arquivos no arquivo de saída.
    echo "" >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"
done

echo "-------------------------------------"
echo "Script concluído com sucesso!"
echo "A estrutura e o conteúdo dos arquivos foram salvos em: $OUTPUT_FILE"