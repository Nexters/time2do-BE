INSERT INTO T2D_DB.participants
VALUES (1, 1);

INSERT INTO T2D_DB.users
VALUES (1, '사람', '123456789');

INSERT INTO T2D_DB.timers
VALUES (1, '개인 타이머', 1, 2, '태그1, 태그2, 태그3', null, '2023-02-01', '2023-02-02');

INSERT INTO T2D_DB.timers
VALUES (2, '그룹 타이머', 1, 2, '태그1, 태그2, 태그3', null, '2023-02-19', '2023-02-20');

INSERT INTO T2D_DB.to_dos
VALUES (1, 1, '밥먹기', false, false, '2023-02-01', null, null, null);

INSERT INTO T2D_DB.to_dos
VALUES (2, 1, '잠자기', false, false, '2023-02-01', null, null, null);

INSERT INTO T2D_DB.to_dos
VALUES (3, 1, '숨쉬기', false, false, '2023-02-01', null, null, null);

INSERT INTO T2D_DB.time_records
VALUES (1, 1, 1, '2023-02-01 00:00:00', '2023-02-01 06:00:00');