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
    id: 1234    // uint
    username: "투두두01"    // string : 닉네임
    password: "passw"    // string : 비밀번호
}
```

<br/>

**<u>group</u>**
```bash
{
    "id": 1,    // uint
    "name": "1번다운타이머",    // string : 타이머 이름
    "makerId": 1234,    // uint : 제작자의 id
    "Participants": "1234.1235.1236",   // string : 참여자들의 ids
    "tags": "열정.개발.디자인", // string : 해시태그
    "setTime": "0001-01-01T00:00:00Z"   // string(datetime) 
}
```

<br/>

**<u>todo</u>**
```bash
{
    "id": 1,    // uint
    "userId": 1234,    // uid : todo 소유자 id
}
```