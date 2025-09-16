# VANTUN - Protocolo de Túnel Seguro de Próxima Geração

VANTUN é um protocolo de túnel de ponta e alto desempenho construído sobre QUIC, projetado para oferecer desempenho de rede excepcional, segurança e confiabilidade. Como uma solução de próxima geração, VANTUN redefine o que é possível no túnel de redes com sua arquitetura inovadora e recursos avançados.

## Principais Vantagens

### 🔒 Segurança Empresarial
- **Handshake Seguro e Negociação de Sessão**: Realizado via stream de controle dedicado para segurança da conexão

### ⚡ Desempenho Excepcional
- **Múltiplos Tipos de Streams Lógicos**: Streams interativos, de grande volume e de telemetria otimizados para diferentes cenários de negócios
- **Multipath**: Uso inteligente de múltiplos caminhos de rede para velocidade e estabilidade de conexão drasticamente melhoradas

### 🛡️ Confiabilidade Incomparável
- **Correção de Erros Progressiva (FEC)**: Tecnologia avançada de correção de erros garante integridade de dados mesmo em condições de rede instáveis
- **Controle de Congestão Híbrido**: Algoritmo híbrido inovador combinando QUIC CC com limitação de taxa por token bucket para utilização ótima de recursos

### 🌐 Proteção de Privacidade
- **Módulo de Ofuscação Plugável**: Ofuscação de tráfego avançada que faz o tráfego parecer HTTP/3 normal, evitando efetivamente o escrutínio da rede

### 🚀 Implantação Fácil
- **Cliente/Servidor Mínimos**: Programas de linha de comando `client` e `server` para implantação rápida e facilidade de uso

## Arquitetura Tecnológica

VANTUN aproveita tecnologias líderes do setor para oferecer seu desempenho e confiabilidade excepcionais:

- **Linguagem**: Go - Linguagem de programação moderna de alto desempenho e concorrente
- **Biblioteca Core**: `quic-go` - Implementação de protocolo QUIC líder no setor
- **Serialização**: `github.com/fxamacker/cbor` - Codificação CBOR eficiente, mais compacta que JSON
- **FEC**: `github.com/klauspost/reedsolomon` - Algoritmo de codificação Reed-Solomon de alto desempenho
- **CLI**: `cobra/viper` - Interface de linha de comando poderosa e gerenciamento de configuração

## Início Rápido

Coloque VANTUN em funcionamento em apenas alguns minutos:

1. **Clonar Repositório**: `git clone <repository-url>`
2. **Construir**: `go build -o bin/vantun cmd/main.go`
3. **Configurar**: Criar arquivo de configuração `config.json`
4. **Executar**: Iniciar servidor e cliente

Para instruções detalhadas e configuração, consulte o [Guia de Demonstração](DEMOGUIDE_pt.md).

## Estrutura do Projeto

```
vantun/
├── cmd/              # Entrada do programa de linha de comando
├── internal/
│   ├── cli/          # Gerenciamento de configuração CLI
│   └── core/         # Implementação do protocolo core
├── docs/             # Documentação
├── go.mod            # Definição do módulo Go
└── README.md         # Documentação do projeto
```

## Destaques da Arquitetura

### 🔧 Motor de Protocolo Inteligente
O motor de protocolo core implementa negociação de sessão eficiente e gerenciamento de stream de controle para conexões seguras e estáveis.

### 📊 Tecnologia FEC Adaptativa
Correção de erros progressiva baseada em codificação Reed-Solomon que ajusta dinamicamente estratégias de correção baseadas em condições de rede.

### 🔄 Transmissão Multipath Inteligente
Gerenciamento de caminho inovador e balanceamento de carga que utiliza completamente todos os caminhos de rede disponíveis para redundância e throughput aprimorado.

### 📈 Controle de Congestão Híbrido
Algoritmo híbrido combinando controle de congestão QUIC subjacente com token bucket de camada superior para utilização ótima de recursos.

### 🎭 Ofuscação de Tráfego Avançada
Ofuscação de tráfego estilo HTTP/3 e preenchimento de dados inteligente para evitar efetivamente o escrutínio da rede e proteger a privacidade do usuário.

### 📊 Sistema de Telemetria em Tempo Real
Coleta abrangente de dados de desempenho e monitoramento em tempo real para otimização de rede e solução de problemas.

## Garantia de Qualidade

VANTUN adota padrões de teste rigorosos para garantir qualidade de código e estabilidade do sistema:

- **Testes Unitários Abrangentes**: Abrangendo todos os módulos funcionais core
- **Testes de Integração**: Validando colaboração de componentes
- **Testes de Desempenho**: Garantindo desempenho excepcional sob várias condições de rede
- **Testes de Estresse**: Validando estabilidade sob alta carga

Executar todos os testes:

```bash
go test -v ./internal/core/...
```

## Licença

VANTUN é licenciado sob a Licença MIT, uma licença open-source permissiva que permite uso, cópia, modificação e distribuição gratuita do software enquanto retém os avisos de direitos autorais e licença.

---

*© 2025 Projeto VANTUN. Todos os direitos reservados.*