# VANTUN - Protocolo de Túnel Seguro de Próxima Generación

VANTUN es un protocolo de túnel de vanguardia y alto rendimiento construido sobre QUIC, diseñado para ofrecer un rendimiento de red excepcional, seguridad y confiabilidad. Como una solución de próxima generación, VANTUN redefine lo que es posible en el túneleo de redes con su arquitectura innovadora y características avanzadas.

## Ventajas Principales

### 🔒 Seguridad de Grado Empresarial
- **Handshake Seguro y Negociación de Sesión**: Realizado a través de un stream de control dedicado para la seguridad de la conexión

### ⚡ Rendimiento Excepcional
- **Múltiples Tipos de Streams Lógicos**: Streams interactivos, de gran volumen y de telemetría optimizados para diferentes escenarios de negocio
- **Multipath**: Uso inteligente de múltiples rutas de red para una velocidad y estabilidad de conexión drásticamente mejoradas

### 🛡️ Confiabilidad Inigualable
- **Corrección de Errores hacia Adelante (FEC)**: Tecnología de corrección de errores avanzada que garantiza la integridad de datos incluso en condiciones de red inestables
- **Control de Congestión Híbrido**: Algoritmo híbrido innovador que combina QUIC CC con limitación de tasa por cubo de tokens para una utilización óptima de recursos

### 🌐 Protección de Privacidad
- **Módulo de Ofuscación Conectable**: Ofuscación de tráfico avanzada que hace que el tráfico parezca HTTP/3 normal, evitando efectivamente el escrutinio de red

### 🚀 Despliegue Fácil
- **Cliente/Servidor Mínimos**: Programas de línea de comando `client` y `server` para despliegue rápido y facilidad de uso

## Arquitectura Tecnológica

VANTUN aprovecha tecnologías líderes en la industria para ofrecer su rendimiento y confiabilidad excepcionales:

- **Lenguaje**: Go - Lenguaje de programación moderno de alto rendimiento y concurrente
- **Biblioteca Core**: `quic-go` - Implementación de protocolo QUIC líder en la industria
- **Serialización**: `github.com/fxamacker/cbor` - Codificación CBOR eficiente, más compacta que JSON
- **FEC**: `github.com/klauspost/reedsolomon` - Algoritmo de codificación Reed-Solomon de alto rendimiento
- **CLI**: `cobra/viper` - Interfaz de línea de comando potente y gestión de configuración

## Inicio Rápido

Ponga VANTUN en funcionamiento en solo unos minutos:

1. **Clonar Repositorio**: `git clone <repository-url>`
2. **Construir**: `go build -o bin/vantun cmd/main.go`
3. **Configurar**: Crear archivo de configuración `config.json`
4. **Ejecutar**: Iniciar servidor y cliente

Para instrucciones detalladas y configuración, por favor refiérase a la [Guía de Demostración](DEMOGUIDE_es.md).

## Estructura del Proyecto

```
vantun/
├── cmd/              # Entrada del programa de línea de comando
├── internal/
│   ├── cli/          # Gestión de configuración CLI
│   └── core/         # Implementación del protocolo core
├── docs/             # Documentación
├── go.mod            # Definición del módulo Go
└── README.md         # Documentación del proyecto
```

## Aspectos Destacados de la Arquitectura

### 🔧 Motor de Protocolo Inteligente
El motor de protocolo core implementa negociación de sesión eficiente y gestión de stream de control para conexiones seguras y estables.

### 📊 Tecnología FEC Adaptativa
Corrección de errores hacia adelante basada en codificación Reed-Solomon que ajusta dinámicamente estrategias de corrección basadas en condiciones de red.

### 🔄 Transmisión Multipath Inteligente
Gestión de rutas innovadora y balanceo de carga que utiliza completamente todas las rutas de red disponibles para redundancia y mayor rendimiento.

### 📈 Control de Congestión Híbrido
Algoritmo híbrido que combina el control de congestión QUIC subyacente con el cubo de tokens de capa superior para una utilización óptima de recursos.

### 🎭 Ofuscación de Tráfico Avanzada
Ofuscación de tráfico de estilo HTTP/3 y relleno de datos inteligente para evitar efectivamente el escrutinio de red y proteger la privacidad del usuario.

### 📊 Sistema de Telemetría en Tiempo Real
Recolección de datos de rendimiento integral y monitoreo en tiempo real para optimización de red y solución de problemas.

## Garantía de Calidad

VANTUN adopta estándares de prueba estrictos para garantizar la calidad del código y la estabilidad del sistema:

- **Pruebas Unitarias Integrales**: Cubriendo todos los módulos funcionales core
- **Pruebas de Integración**: Validando la colaboración de componentes
- **Pruebas de Rendimiento**: Asegurando rendimiento excepcional bajo diversas condiciones de red
- **Pruebas de Estrés**: Validando estabilidad bajo alta carga

Ejecutar todas las pruebas:

```bash
go test -v ./internal/core/...
```

## Licencia

VANTUN está licenciado bajo la Licencia MIT, una licencia de código abierto permisiva que permite el uso, copia, modificación y distribución gratuita del software mientras se retienen los avisos de derechos de autor y licencia.

---

*© 2025 Proyecto VANTUN. Todos los derechos reservados.*