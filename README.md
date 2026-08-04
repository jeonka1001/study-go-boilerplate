# Study Go Boilerplate

`study-go-boilerplate`는 Go 언어를 사용하여 확장 가능하고 모듈화된 RESTful API 서버를 신속하게 구축할 수 있도록 설계된 보일러플레이트(Boilerplate) 프로젝트입니다.

Uber의 **Fx** 프레임워크를 이용한 의존성 주입(DI)과 **Fiber v2** 웹 프레임워크를 기반으로 도메인 중심 디렉토리 구조 및 레이어드 아키텍처를 제공합니다.

---

## 🛠 기술 스택 (Tech Stack)

| 구분 | 라이브러리 / 도구 | 설명 |
| :--- | :--- | :--- |
| **Language** | [Go 1.24+](https://go.dev/) | 백엔드 개발 언어 |
| **Web Framework** | [Fiber v2](https://github.com/gofiber/fiber) | Express 스타일의 고성능 Go 웹 프레임워크 |
| **DI Framework** | [Uber Fx](https://github.com/uber-go/fx) | 의존성 주입(Dependency Injection) 및 생명주기 관리 |
| **ORM** | [GORM](https://gorm.io/) | MySQL 등의 RDB 연동 및 데이터 매핑 |
| **Configuration** | [Viper](https://github.com/spf13/viper) | YAML 기반 환경 설정 관리 |
| **Logger** | [Uber Zap](https://github.com/uber-go/zap) | 고성능 구조화 로깅 |
| **Validation** | [Validator v10](https://github.com/go-playground/validator) | DTO 및 요청 데이터 검증 |
| **API Docs** | [Swagger](https://github.com/gofiber/swagger) | API 문서화 |

---

## 📁 폴더 구조 및 구성 요소 상세 (Detailed Folder Structure)

```text
study-go-boilerplate/
├── main.go                       # 서버 엔트리포인트 (Uber Fx 의존성 그래프 조립 및 실행)
├── go.mod                        # Go 모듈 및 의존성 라이브러리 정의
├── go.sum                        # 의존성 패키지 체크섬
│
├── config/                       # 환경 설정(Configuration) 관리 모듈
│   ├── config.go                 # APP_ENV 환경변수를 파싱하여 해당 profile의 yml 파일 및 시스템 환경변수 오버라이드 로드
│   ├── module.go                 # config 패키지의 Fx DI 모듈 (fx.Provide)
│   ├── type.go                   # Server, App, Log, DB 관련 설정 구조체(Struct) 정의
│   └── yml/
│       └── local.yml             # 로컬 개발 환경용 서버/DB/로거 기본 설정값
│
└── app/                          # 애플리케이션 핵심 서비스 및 비즈니스 도메인
    ├── app.go                    # Fiber 객체 초기화, Lifecycle(OnStart/OnStop) 바인딩, Swagger 및 Health Check 라우트 등록
    ├── middleware.go             # 전역 미들웨어(Recover, CORS, RequestID, Zap HTTP Logger) 순서 지정 및 등록
    ├── router.go                 # Fx Grouping(`group:"routes"`)을 활용하여 등록된 모든 컨트롤러의 라우트 매핑 테이블을 Fiber 앱에 자동 바인딩
    │
    ├── core/                     # 공통 아키텍처 인프라, 유틸리티, 헬퍼
    │   ├── module.go             # Core 관련 Fx 모듈 통합 정의 (`BaseModule`, `RepositoryModule` 등)
    │   ├── base/                 # 공통 요청/응답 파라미터 파서 및 바인더
    │   │   └── parameter.go      # Path, Header, Query, Body의 자동 파싱 및 validator 기반 DTO 유효성 검증 인터페이스(`Parameter`)
    │   ├── client/               # 외부 서비스 호출용 HTTP/RPC 클라이언트 모듈 위치
    │   │   └── client.go
    │   ├── consts/               # 프로젝트 전역 상수 정의
    │   │   └── const.go
    │   ├── dao/                  # Data Access Object 및 데이터베이스 Connection 관리
    │   │   ├── mysql.go          # GORM 기반 MySQL 연결 설정, Connection Pool 옵션 적용 및 OnStart/OnStop 생명주기 관리
    │   │   └── kafka.go          # 카프카 데이터 액세스 객체 모듈 위치
    │   ├── exception/            # 예외 처리 및 커스텀 에러 처리
    │   │   ├── error.go          # 시스템 에러 코드 정의 (BadRequest, Unauthorized, InternalServerError 등) 및 커스텀 Error 객체
    │   │   └── func.go           # Fiber 전역 ErrorHandler 미들웨어 및 Panic Recover 미들웨어 구현
    │   ├── helper/               # 로거 및 유효성 검증기 모듈
    │   │   ├── module.go         # Helper Fx 모듈
    │   │   ├── logger/           # Uber Zap 기반 로깅 모듈 (zap.go, logger.go, config.go)
    │   │   └── validator/        # validator/v10 기반 구조체 유효성 검증 객체 생성
    │   ├── repository/           # 데이터 액세스 리포지토리 레이어
    │   │   └── repository.go     # GORM DB 인스턴스를 주입받는 공통 Repository 래퍼 구조체
    │   └── util/                 # 공통 유틸리티 함수
    │       └── time.go           # 날짜 및 시간 처리 유틸리티
    │
    └── domain/                   # 도메인 별 비즈니스 로직 (DDD 스타일 모듈화)
        └── user/                 # 'User' 도메인 예시 모듈
            ├── module.go         # User 도메인의 ControllerModule, ServiceModule 정의 (AsRoute를 통한 라우트 자동 등록)
            ├── controller/       # HTTP 요청 처리기 (`Table() []app.Mapping` 인터페이스 구현을 통해 엔드포인트 정의)
            │   └── controller.go
            ├── dto/              # API 요청/응답 DTO (Data Transfer Object)
            │   ├── request.go    # 사용자 생성/수정 등 요청 Body 및 Param 구조체
            │   └── resposne.go   # 사용자 정보 응답 구조체
            └── service/          # 비즈니스 핵심 로직 처리 계층
                └── service.go    # User 서비스 인터페이스 및 구현체 (Repository 참조)
```

---

## 🔍 디렉토리별 역할 상세 설명

### 1. `config/`
- **역할**: 애플리케이션의 동작 환경(Local, Dev, QA, Prod)에 따른 모든 설정값을 관리합니다.
- **주요 파일**:
  - `config.go`: `APP_ENV` 환경변수(기본값 `local`)를 확인하고 `config/yml/{APP_ENV}.yml` 경로의 설정 파일을 Viper로 읽습니다. 또한 `DB_PASSWORD` 등 환경변수를 통한 설정값 오버라이딩을 지원합니다.
  - `type.go`: YAML 설정과 1:1로 매핑되는 `Config`, `Server`, `DB`, `Log` 등 Go 구조체를 다룹니다.

### 2. `app/` (Core App Module)
- **역할**: 웹 프레임워크(Fiber)의 수명주기, 글로벌 미들웨어, 라우팅 설정을 전담합니다.
- **주요 파일**:
  - `app.go`: Fiber 앱 인스턴스(`initializeFiber`)를 생성하고, Uber Fx `Lifecycle`에 서버의 `Listen` 및 `Shutdown`을 바인딩합니다. `/check_health` 및 Swagger UI 라우트도 이곳에 등록됩니다.
  - `router.go`: `AsRoute` 헬퍼 함수로 태깅된 도메인 컨트롤러들을 전달받아, 각 컨트롤러가 가진 라우트 테이블(`Table() []Mapping`)을 Fiber 앱 라우터에 자동으로 바인딩합니다.
  - `middleware.go`: Recover, CORS, RequestID, Zap HTTP Logger 전역 미들웨어를 순서대로 체이닝하여 등록합니다.

### 3. `app/core/` (Common Infrastructure Layer)
- **역할**: 도메인 영역에 종속되지 않고 프로젝트 전반에서 재사용되는 인프라, DB, 로깅, 에러 처리, 파라미터 파서 등을 제공합니다.
- **주요 패키지**:
  - `base/parameter.go`: Client의 HTTP Request(Query, Path, Body, Header)를 DTO로 파싱하고 `validator`를 통한 검증 과정을 캡슐화한 `Parameter` 인터페이스를 제공합니다.
  - `dao/mysql.go`: GORM을 사용해 MySQL에 접속하고 Connection Pool 설정(MaxOpen, MaxIdle, Lifetime 등)을 적용합니다. Fx `Lifecycle`에 접속 테스트(Ping) 및 Termination 시 Connection Close 동작이 포함되어 있습니다.
  - `exception/`: 서비스 전반의 공통 에러 코드(`CodeBadRequest`, `CodeNotFound` 등)와 `CustomError` 타입을 정의합니다. `ErrorHandler`는 발생한 에러를 일관된 JSON 형식으로 Client에 반환합니다.
  - `helper/logger/`: Uber Zap 로거를 설정하고 Fiber 전역 미들웨어로 연동하여 요청 별 ID, Latency, Status, IP 등을 구조화된 로그로 남깁니다.

### 4. `app/domain/` (Business Domain Layer)
- **역할**: 실제 도메인 특화 비즈니스 로직이 구현되는 영역입니다. 도메인 단위(예: `user`, `order`, `product` 등)로 폴더가 격리되어 관리됩니다.
- **도메인 구성요소 (예: `user`)**:
  - `controller/`: HTTP 요청 수신, `Parameter.GetRequest` 및 `ValidateParams`를 이용한 요청 검증, `service` 호출 후 응답 반환을 담당합니다. `app.Route` 인터페이스를 구현하여 자신의 라우트 경로와 HTTP 메소드 매핑을 제공합니다.
  - `service/`: 비즈니스 유효성 검사 및 데이터 처리를 담당합니다. 필요 시 `core.Repository`를 주입받아 데이터베이스 조회를 수행합니다.
  - `dto/`: Request / Response 전용 구조체를 정의합니다.
  - `module.go`: Fx 모듈 단위로 Controller와 Service의 생성자를 등록하며, `app.AsRoute`를 통해 컨트롤러를 Fx 라우터 그룹에 결합시킵니다.

---

## 🚀 주요 핵심 특징 (Key Highlights)

1. **Uber Fx 기반의 의존성 주입 (DI)**
   - 모듈(`fx.Module`) 단위로 의존성을 정의하며, `main.go`에서 통합 관리합니다.
   - 새로운 서비스나 컨트롤러 생성 시 Fx 모듈에 등록하기만 하면 자동으로 주입됩니다.

2. **자동 라우트 등록 (Annotated Route Grouping)**
   - `app/router.go`의 `AsRoute` 헬퍼 함수를 사용하여 컨트롤러를 Fx의 `"routes"` 그룹으로 바인딩합니다.
   - 새로운 컨트롤러 추가 시 라우트를 수동으로 개별 등록할 필요 없이 `Route` 인터페이스만 구현하면 자동 연동됩니다.

3. **중앙 집중식 커스텀 에러 핸들링**
   - `app/core/exception` 및 `app/middleware.go`에서 예외를 세분화하여 공통 에러 응답 형식으로 변환 처리합니다.

4. **환경별 설정 관리 (Viper & YAML)**
   - `config/yml/` 경로에 실행 환경에 따른 YAML 파일을 배치하여 체계적으로 환경 변수를 관리합니다.

---

## 🏁 시작하기 (Getting Started)

### 1. 사전 요구사항 (Prerequisites)
- [Go 1.24](https://go.dev/doc/install) 이상 설치 필요

### 2. 프로젝트 모듈명 변경 (선택 사항)
본 프로젝트의 기본 모듈명은 `boilerplate`로 되어 있습니다. 실제 사용하려는 저장소 주소로 변경하려면 `go.mod` 및 프로젝트 내 import 구문을 수정하세요.
```bash
# go.mod 모듈명 변경 예시
go mod edit -module github.com/your-username/your-repo
```

### 3. 의존성 패키지 설치
```bash
go mod download
```

### 4. 서버 실행 (Run)
```bash
go run main.go
```

---

## 💡 새로운 도메인 기능 추가 가이드

1. `app/domain/{new_domain}` 디렉토리를 생성합니다.
2. 하위에 `controller/`, `service/`, `dto/` 및 `module.go`를 작성합니다.
3. `controller` 구현체에 `app.Route` 인터페이스(`Table() []app.Mapping`)를 구현합니다.
4. `module.go`에서 `app.AsRoute(NewController)`로 라우트를 Fx 모듈에 등록합니다.
5. `main.go`의 `fx.New()` 목록에 신규 도메인의 모듈을 추가합니다.
