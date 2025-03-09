### 🌟 프로젝트 기술 설명서 🌟

이 프로젝트는 **파일 호스팅**를 위한 **RESTful API 서버**입니다. Go 언어로 작성되었으며, Docker를 이용해 배포할 수 있도록 구성되었습니다. 🚀 서버는 사용자가 인증된 상태에서 파일을 업로드하고 관리할 수 있도록 기능을 제공합니다. 

---

### 🧑‍💻 주요 기능 🧑‍💻

1. **파일 업로드**: 인증된 사용자만 파일을 업로드할 수 있습니다. 파일은 사용자별로 구분된 폴더에 저장됩니다. 📂
2. **파일 수정**: 인증된 사용자는 자신이 업로드한 파일을 수정할 수 있습니다. (파일명 및 내용 수정 가능) ✏️
3. **파일 다운로드**: 모든 사용자는 인증 없이 모든 파일을 다운로드할 수 있습니다. ⬇️
4. **파일 삭제**: 인증된 사용자는 자신이 업로드한 파일을 삭제할 수 있습니다. 🗑️
5. **사용자 관리**: 관리자는 새로운 사용자를 등록할 수 있으며, 각 사용자에게 자동으로 고유한 비밀번호를 생성하여 해시값을 저장합니다. 🔐

---

### 🛠️ 기술 스택 🛠️

- **Go**: 서버 사이드 로직 구현
- **Docker**: 애플리케이션 컨테이너화 (개발 및 프로덕션 환경 지원)
- **Gorilla Mux**: HTTP 라우팅
- **bcrypt**: 안전한 비밀번호 해싱
- **JSON**: 데이터 저장 및 응답 형식으로 사용

---

### 📁 프로젝트 구조 📁

```
📁 img-host-server
│
├── /cmd/server/main.go          : 서버 초기화 및 라우팅 설정
│
├── /internal
│   ├── handles
│   │   ├── file_handler.go      : 파일 업로드, 수정, 삭제, 다운로드 처리
│   │   └── user_handler.go      : 사용자 등록 및 관리 처리
│   ├── utils
│   │   ├── auth.go              : 사용자 인증 관련 함수 (비밀번호 검증, 사용자 정보 로딩)
│   │   ├── response.go          : JSON 형식의 응답을 처리하는 유틸리티
│   │   ├── file.go              : 파일 업로드 및 저장 관련 함수
│   │   └── sanitize.go          : 파일명 유효성 검사
│   └── db
│       └── users.json           : 사용자 정보 저장
│
├── .air.toml                    : 개발 환경 hotswap 라이브러리 설정 파일(air 라이브러리)
├── .env                         : 환경변수 관리(관리자 비밀번호)
├── /docker-compose.yml          : Docker 설정 파일
├── /Dockerfile.prod             : 프로덕션 환경을 위한 Dockerfile
├── /Dockerfile.dev              : 개발 환경을 위한 Dockerfile
├── /go.mod                      : Go 모듈 설정 파일
└── /go.sum                      : Go 의존성 정보 파일

```

---

### 💻 API 엔드포인트 💻

- **POST `/files`**: 파일 업로드
- **PUT `/files/{filename}`**: 파일 수정
- **GET `/files/{username}/{filename}`**: 파일 다운로드
- **DELETE `/files/{filename}`**: 파일 삭제
- **POST `/users`**: 사용자 등록

---

### 🚀 프로젝트 실행 방법 🚀

GitHub에서 이 프로젝트를 **Clone**한 후, 로컬에서 실행하는 방법은 다음과 같습니다. 🔽

#### 1️⃣ **GitHub에서 프로젝트 Clone**
```bash
git clone https://github.com/Aleph-Kim/img-host-server
cd img-host-server
```

#### 2️⃣ **관리자 비밀번호 설정**
`.env` 파일을 사용하여 관리자 비밀번호를 설정할 수 있습니다.

#### 3️⃣ **Docker 환경 설정**
Docker를 이용하여 애플리케이션을 실행할 수 있습니다. `docker-compose.yml` 파일을 사용하여 개발 환경과 프로덕션 환경을 설정할 수 있습니다.

- **개발 환경 실행**:
```bash
docker-compose up dev
```

- **프로덕션 환경 실행**:
```bash
docker-compose up prod
```

#### 4️⃣ **서버 실행**
서버는 기본적으로 `3000` 포트에서 실행됩니다. 서버가 성공적으로 실행되면, 브라우저나 Postman 등을 통해 API를 테스트할 수 있습니다.

---

### 📦 Docker 구성 설명 📦

- `Dockerfile.prod`: 프로덕션 환경용 Dockerfile로, 실제 배포용 설정이 포함되어 있습니다.
- `Dockerfile.dev`: 개발 환경용 Dockerfile로, 소스 코드 변경 시 자동 반영되도록 설정되어 있습니다.
- `docker-compose.yml`: 개발 및 프로덕션 환경을 동시에 지원하는 설정 파일로, 필요한 서비스를 자동으로 구성해줍니다.

---

### 🚨 주의 사항 🚨

- `.env` 파일에 관리자의 비밀번호가 설정되어야 하며, 관리자 인증을 통해 사용자를 등록할 수 있습니다.
- 파일 업로드/수정/삭제 시 사용자 인증을 위해 **X-Username**과 **X-Secret** 헤더를 사용합니다.

---

이 프로젝트는 파일 업로드와 관리, 사용자 인증을 간편하고 안전하게 처리할 수 있도록 설계되었습니다. 추가적인 기능이 필요하면 언제든지 요구 사항을 반영할 수 있습니다. 😊
