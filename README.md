# Time2Do Backend

<br/>

### Command

```bash
# 누락된 모듈 추가 & 사용하지 않는 모듈 제거
$ go mod tidy

# go.sum 파일 유효성 검사
$ go mod verify

# 서버 실행
bash run.sh
```

<br/>
<br/>

### Data

<br/>

**<u>user</u>**

```bash
{
    id: 1234    // int
    username: "투두두01"    // string 닉네임
    password: "passw"    // string 비밀번호
}
```

<br/>

**<u>group</u>**
```bash
{
    "id": 1,
    "name": "1번다운타이머",
    "makerId": 1234,
    "Participants": "1234.1235.1236",
    "tags": "열정.개발.디자인",
    "setTime": "0001-01-01T00:00:00Z"
}
```