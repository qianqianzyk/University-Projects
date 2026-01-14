-- 1. 表结构设计
-- 1.1 专业表 department
CREATE TABLE Zhaoyk_department (
    zyk_id SERIAL PRIMARY KEY,
    zyk_name VARCHAR(100) NOT NULL UNIQUE,
    zyk_create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    zyk_update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE Zhaoyk_department IS '专业信息表';
COMMENT ON COLUMN Zhaoyk_department.zyk_id IS '专业编号';
COMMENT ON COLUMN Zhaoyk_department.zyk_name IS '专业名称';
COMMENT ON COLUMN Zhaoyk_department.zyk_create_time IS '创建时间';
COMMENT ON COLUMN Zhaoyk_department.zyk_update_time IS '更新时间';

-- 1.2 班级表 class
CREATE TABLE Zhaoyk_class (
    zyk_id SERIAL PRIMARY KEY,
    zyk_name VARCHAR(100) NOT NULL,
    zyk_department_id INT NOT NULL,
    zyk_create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    zyk_update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE Zhaoyk_class IS '班级信息表';
COMMENT ON COLUMN Zhaoyk_class.zyk_id IS '班级编号';
COMMENT ON COLUMN Zhaoyk_class.zyk_name IS '班级名称';
COMMENT ON COLUMN Zhaoyk_class.zyk_department_id IS '所属专业编号';
COMMENT ON COLUMN Zhaoyk_class.zyk_create_time IS '创建时间';
COMMENT ON COLUMN Zhaoyk_class.zyk_update_time IS '更新时间';

CREATE INDEX Zhaoyk_idx_class_department ON Zhaoyk_class(zyk_department_id);
ALTER TABLE Zhaoyk_class ADD CONSTRAINT Zhaoyk_fk_class_department FOREIGN KEY (zyk_department_id) REFERENCES Zhaoyk_department(zyk_id) ON DELETE RESTRICT;

-- 1.3 学生表 student
CREATE TABLE Zhaoyk_stu (
    zyk_id SERIAL PRIMARY KEY,
    zyk_student_id CHAR(12) NOT NULL UNIQUE,
    zyk_password VARCHAR(128) NOT NULL,
    zyk_name VARCHAR(50) NOT NULL,
    zyk_gender CHAR(1),
    zyk_age INT,
    zyk_city_id INT,
    zyk_class_id INT NOT NULL,
    zyk_gpa NUMERIC(4,2) DEFAULT 0,
    zyk_credits NUMERIC(5,2) DEFAULT 0,
    zyk_create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    zyk_update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE Zhaoyk_stu IS '学生信息表';
COMMENT ON COLUMN Zhaoyk_stu.zyk_id IS '学生编号';
COMMENT ON COLUMN Zhaoyk_stu.zyk_student_id IS '学号';
COMMENT ON COLUMN Zhaoyk_stu.zyk_password IS '密码';
COMMENT ON COLUMN Zhaoyk_stu.zyk_name IS '姓名';
COMMENT ON COLUMN Zhaoyk_stu.zyk_gender IS '性别（M/F）';
COMMENT ON COLUMN Zhaoyk_stu.zyk_age IS '年龄';
COMMENT ON COLUMN Zhaoyk_stu.zyk_city_id IS '生源地城市编号';
COMMENT ON COLUMN Zhaoyk_stu.zyk_class_id IS '所属班级编号';
COMMENT ON COLUMN Zhaoyk_stu.zyk_gpa IS '学生平均绩点';
COMMENT ON COLUMN Zhaoyk_stu.zyk_credits IS '已修总学分';
COMMENT ON COLUMN Zhaoyk_stu.zyk_create_time IS '创建时间';
COMMENT ON COLUMN Zhaoyk_stu.zyk_update_time IS '更新时间';

CREATE INDEX Zhaoyk_idx_stu_city ON Zhaoyk_stu(zyk_city_id);
CREATE INDEX Zhaoyk_idx_stu_class ON Zhaoyk_stu(zyk_class_id);
ALTER TABLE Zhaoyk_stu ADD CONSTRAINT Zhaoyk_student_zyk_name_check CHECK (TRIM(zyk_name) <> '');
ALTER TABLE Zhaoyk_stu ADD CONSTRAINT Zhaoyk_student_zyk_gender_check CHECK (zyk_gender IN ('M', 'F'));
ALTER TABLE Zhaoyk_stu ADD CONSTRAINT Zhaoyk_student_zyk_age_check CHECK (zyk_age BETWEEN 1 AND 100);
ALTER TABLE Zhaoyk_stu ADD CONSTRAINT Zhaoyk_fk_student_class FOREIGN KEY (zyk_class_id) REFERENCES Zhaoyk_class(zyk_id) ON DELETE RESTRICT;

-- 1.4 教师表 teacher
CREATE TABLE Zhaoyk_tea (
    zyk_id SERIAL PRIMARY KEY,
    zyk_teacher_id VARCHAR(10) NOT NULL UNIQUE,
    zyk_password VARCHAR(128) NOT NULL,
    zyk_name VARCHAR(50) NOT NULL,
    zyk_gender CHAR(1),
    zyk_age INT,
    zyk_title VARCHAR(50),
    zyk_phone VARCHAR(20),
    zyk_is_admin BOOLEAN DEFAULT FALSE,
    zyk_create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    zyk_update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE Zhaoyk_tea IS '教师信息表';
COMMENT ON COLUMN Zhaoyk_tea.zyk_id IS '教师编号';
COMMENT ON COLUMN Zhaoyk_tea.zyk_teacher_id IS '教师工号';
COMMENT ON COLUMN Zhaoyk_tea.zyk_password IS '教师密码';
COMMENT ON COLUMN Zhaoyk_tea.zyk_name IS '姓名';
COMMENT ON COLUMN Zhaoyk_tea.zyk_gender IS '性别';
COMMENT ON COLUMN Zhaoyk_tea.zyk_age IS '年龄';
COMMENT ON COLUMN Zhaoyk_tea.zyk_title IS '职称';
COMMENT ON COLUMN Zhaoyk_tea.zyk_phone IS '联系电话';
COMMENT ON COLUMN Zhaoyk_tea.zyk_is_admin IS '是否为管理员（TRUE 表示是管理员）';
COMMENT ON COLUMN Zhaoyk_tea.zyk_create_time IS '创建时间';
COMMENT ON COLUMN Zhaoyk_tea.zyk_update_time IS '更新时间';

ALTER TABLE Zhaoyk_tea ADD CONSTRAINT Zhaoyk_tea_zyk_name_check CHECK (TRIM(zyk_name) <> '');
ALTER TABLE Zhaoyk_tea ADD CONSTRAINT Zhaoyk_tea_zyk_gender_check CHECK (zyk_gender IN ('M', 'F'));
ALTER TABLE Zhaoyk_tea ADD CONSTRAINT Zhaoyk_tea_zyk_age_check CHECK (zyk_age BETWEEN 1 AND 100);
ALTER TABLE Zhaoyk_tea ADD CONSTRAINT Zhaoyk_tea_zyk_phone_check CHECK (zyk_phone ~ '^1\d{10}$' OR zyk_phone ~ '^\d{3,4}-?\d{7,8}$');

-- 1.5 课程表 course
CREATE TABLE Zhaoyk_cou (
    zyk_id SERIAL PRIMARY KEY,
    zyk_name VARCHAR(100) NOT NULL,
    zyk_school_year INT,
    zyk_semester CHAR(1),
    zyk_hours INT,
    zyk_credit NUMERIC(4,2),
    zyk_class_id INT NOT NULL,
    zyk_exam_type VARCHAR(10),
    zyk_create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    zyk_update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE Zhaoyk_cou IS '课程信息表';
COMMENT ON COLUMN Zhaoyk_cou.zyk_id IS '课程编号';
COMMENT ON COLUMN Zhaoyk_cou.zyk_name IS '课程名称';
COMMENT ON COLUMN Zhaoyk_cou.zyk_school_year IS '开课年份';
COMMENT ON COLUMN Zhaoyk_cou.zyk_semester IS '开课学期';
COMMENT ON COLUMN Zhaoyk_cou.zyk_hours IS '学时';
COMMENT ON COLUMN Zhaoyk_cou.zyk_credit IS '学分';
COMMENT ON COLUMN Zhaoyk_cou.zyk_class_id IS '归属班级';
COMMENT ON COLUMN Zhaoyk_cou.zyk_exam_type IS '考核方式（考试/考查）';
COMMENT ON COLUMN Zhaoyk_cou.zyk_create_time IS '创建时间';
COMMENT ON COLUMN Zhaoyk_cou.zyk_update_time IS '更新时间';

CREATE INDEX Zhaoyk_idx_cou_year ON Zhaoyk_cou(zyk_school_year);
ALTER TABLE Zhaoyk_cou ADD CONSTRAINT Zhaoyk_cou_zyk_school_year_check CHECK (zyk_school_year >= 2000 AND zyk_school_year <= 2100);
ALTER TABLE Zhaoyk_cou ADD CONSTRAINT Zhaoyk_cou_zyk_semester_check CHECK (zyk_semester IN ('1', '2'));
ALTER TABLE Zhaoyk_cou ADD CONSTRAINT Zhaoyk_cou_zyk_hours_check CHECK (zyk_hours > 0);
ALTER TABLE Zhaoyk_cou ADD CONSTRAINT Zhaoyk_cou_zyk_credit_check CHECK (zyk_credit >= 1 AND zyk_credit <= 4);
ALTER TABLE Zhaoyk_cou ADD CONSTRAINT Zhaoyk_fk_cou_class FOREIGN KEY (zyk_class_id) REFERENCES Zhaoyk_class(zyk_id) ON DELETE RESTRICT;
ALTER TABLE Zhaoyk_cou ADD CONSTRAINT Zhaoyk_cou_zyk_exam_type_check CHECK (zyk_exam_type IN ('考试', '考查'));

-- 1.6 成绩表 score
CREATE TABLE Zhaoyk_score (
    zyk_id SERIAL PRIMARY KEY,
    zyk_student_id CHAR(12) NOT NULL,
    zyk_course_id INT NOT NULL,
    zyk_score NUMERIC(5,2),
    zyk_create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    zyk_update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (zyk_student_id, zyk_course_id)
);

COMMENT ON TABLE Zhaoyk_score IS '学生成绩表';
COMMENT ON COLUMN Zhaoyk_score.zyk_id IS '成绩记录编号';
COMMENT ON COLUMN Zhaoyk_score.zyk_student_id IS '学生学号';
COMMENT ON COLUMN Zhaoyk_score.zyk_course_id IS '课程编号';
COMMENT ON COLUMN Zhaoyk_score.zyk_score IS '课程分数';
COMMENT ON COLUMN Zhaoyk_score.zyk_create_time IS '创建时间';
COMMENT ON COLUMN Zhaoyk_score.zyk_update_time IS '更新时间';

CREATE INDEX Zhaoyk_idx_score_student ON Zhaoyk_score(zyk_student_id);
CREATE INDEX Zhaoyk_idx_score_course ON Zhaoyk_score(zyk_course_id);
ALTER TABLE Zhaoyk_score ADD CONSTRAINT Zhaoyk_score_zyk_score_check CHECK (zyk_score >= 0 AND zyk_score <= 100);
ALTER TABLE Zhaoyk_score ADD CONSTRAINT Zhaoyk_fk_score_student FOREIGN KEY (zyk_student_id) REFERENCES Zhaoyk_stu(zyk_student_id) ON DELETE CASCADE;
ALTER TABLE Zhaoyk_score ADD CONSTRAINT Zhaoyk_fk_score_course FOREIGN KEY (zyk_course_id) REFERENCES Zhaoyk_cou(zyk_id) ON DELETE CASCADE;

-- 1.7 授课表 teaching
CREATE TABLE Zhaoyk_teaching (
    zyk_id SERIAL PRIMARY KEY,
    zyk_teacher_id VARCHAR(10) NOT NULL,
    zyk_course_id INT NOT NULL,
    zyk_assign_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    zyk_create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    zyk_update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (zyk_teacher_id, zyk_course_id)
);

COMMENT ON TABLE Zhaoyk_teaching IS '授课分配表';
COMMENT ON COLUMN Zhaoyk_teaching.zyk_id IS '记录编号';
COMMENT ON COLUMN Zhaoyk_teaching.zyk_teacher_id IS '教师工号';
COMMENT ON COLUMN Zhaoyk_teaching.zyk_course_id IS '课程编号';
COMMENT ON COLUMN Zhaoyk_teaching.zyk_assign_time IS '授课分配时间';
COMMENT ON COLUMN Zhaoyk_teaching.zyk_create_time IS '创建时间';
COMMENT ON COLUMN Zhaoyk_teaching.zyk_update_time IS '更新时间';

CREATE INDEX Zhaoyk_idx_teaching_teacher ON Zhaoyk_teaching(zyk_teacher_id);
CREATE INDEX Zhaoyk_idx_teaching_course ON Zhaoyk_teaching(zyk_course_id);
ALTER TABLE Zhaoyk_teaching ADD CONSTRAINT Zhaoyk_fk_teaching_teacher FOREIGN KEY (zyk_teacher_id) REFERENCES Zhaoyk_tea(zyk_teacher_id) ON DELETE RESTRICT;
ALTER TABLE Zhaoyk_teaching ADD CONSTRAINT Zhaoyk_fk_teaching_course FOREIGN KEY (zyk_course_id) REFERENCES Zhaoyk_cou(zyk_id) ON DELETE RESTRICT;

-- 1.8 省份表 province
CREATE TABLE Zhaoyk_province (
    zyk_id SERIAL PRIMARY KEY,
    zyk_name VARCHAR(100) NOT NULL UNIQUE,
    zyk_create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    zyk_update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE Zhaoyk_province IS '省份信息表';
COMMENT ON COLUMN Zhaoyk_province.zyk_id IS '省份编号';
COMMENT ON COLUMN Zhaoyk_province.zyk_name IS '省份名称';
COMMENT ON COLUMN Zhaoyk_province.zyk_create_time IS '创建时间';
COMMENT ON COLUMN Zhaoyk_province.zyk_update_time IS '更新时间';

-- 1.9 城市表 city
CREATE TABLE Zhaoyk_city (
    zyk_id SERIAL PRIMARY KEY,
    zyk_name VARCHAR(100) NOT NULL,
    zyk_province_id INT NOT NULL,
    zyk_create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    zyk_update_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE Zhaoyk_city IS '城市信息表';
COMMENT ON COLUMN Zhaoyk_city.zyk_id IS '城市编号';
COMMENT ON COLUMN Zhaoyk_city.zyk_name IS '城市名称';
COMMENT ON COLUMN Zhaoyk_city.zyk_province_id IS '所属省份编号';
COMMENT ON COLUMN Zhaoyk_city.zyk_create_time IS '创建时间';
COMMENT ON COLUMN Zhaoyk_city.zyk_update_time IS '更新时间';

CREATE INDEX Zhaoyk_idx_city_province ON Zhaoyk_city(zyk_province_id);
ALTER TABLE Zhaoyk_city ADD CONSTRAINT Zhaoyk_fk_city_province FOREIGN KEY (zyk_province_id) REFERENCES Zhaoyk_province(zyk_id) ON DELETE RESTRICT;

-- 2. 触发器设计
-- 2.1 时间更新触发器
CREATE OR REPLACE FUNCTION Zhaoyk_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  IF TG_OP = 'INSERT' THEN
    NEW.zyk_create_time := CURRENT_TIMESTAMP;
    NEW.zyk_update_time := CURRENT_TIMESTAMP;
  ELSIF TG_OP = 'UPDATE' THEN
    NEW.zyk_update_time := CURRENT_TIMESTAMP;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- department 表触发器
CREATE TRIGGER Zhaoyk_trg_timestamp_department
BEFORE INSERT OR UPDATE ON Zhaoyk_department
FOR EACH ROW
EXECUTE PROCEDURE Zhaoyk_update_timestamp();

-- class 表触发器
CREATE TRIGGER Zhaoyk_trg_timestamp_class
BEFORE INSERT OR UPDATE ON Zhaoyk_class
FOR EACH ROW
EXECUTE PROCEDURE Zhaoyk_update_timestamp();

-- student 表触发器
CREATE TRIGGER Zhaoyk_trg_timestamp_stu
BEFORE INSERT OR UPDATE ON Zhaoyk_stu
FOR EACH ROW
EXECUTE PROCEDURE Zhaoyk_update_timestamp();

-- teacher 表触发器
CREATE TRIGGER Zhaoyk_trg_timestamp_tea
BEFORE INSERT OR UPDATE ON Zhaoyk_tea
FOR EACH ROW
EXECUTE PROCEDURE Zhaoyk_update_timestamp();

-- course 表触发器
CREATE TRIGGER Zhaoyk_trg_timestamp_cou
BEFORE INSERT OR UPDATE ON Zhaoyk_cou
FOR EACH ROW
EXECUTE PROCEDURE Zhaoyk_update_timestamp();

-- score 表触发器
CREATE TRIGGER Zhaoyk_trg_timestamp_score
BEFORE INSERT OR UPDATE ON Zhaoyk_score
FOR EACH ROW
EXECUTE PROCEDURE Zhaoyk_update_timestamp();

-- province 表触发器
CREATE TRIGGER Zhaoyk_trg_timestamp_province
BEFORE INSERT OR UPDATE ON Zhaoyk_province
FOR EACH ROW
EXECUTE PROCEDURE Zhaoyk_update_timestamp();

-- city 表触发器
CREATE TRIGGER Zhaoyk_trg_timestamp_city
BEFORE INSERT OR UPDATE ON Zhaoyk_city
FOR EACH ROW
EXECUTE PROCEDURE Zhaoyk_update_timestamp();

-- teaching 表触发器
CREATE TRIGGER Zhaoyk_trg_timestamp_teaching
BEFORE INSERT OR UPDATE ON Zhaoyk_teaching
FOR EACH ROW
EXECUTE PROCEDURE Zhaoyk_update_timestamp();

-- 2.2 成绩学分触发器
CREATE OR REPLACE FUNCTION Zhaoyk_update_student_gpa()
    RETURNS TRIGGER AS $$
DECLARE
    zyk_total_credits NUMERIC(10,2);
    zyk_total_grade_points NUMERIC(10,2);
    zyk_calculated_gpa NUMERIC(4,2);
BEGIN
    SELECT
        COALESCE(SUM(CASE WHEN e.zyk_score >= 60 THEN c.zyk_credit ELSE 0 END), 0),
        COALESCE(SUM(
                         CASE
                             WHEN e.zyk_score >= 60 THEN c.zyk_credit * ((e.zyk_score - 50) / 10)
                             ELSE 0
                             END
                 ), 0)
    INTO zyk_total_credits, zyk_total_grade_points
    FROM Zhaoyk_score e
    JOIN Zhaoyk_cou c ON e.zyk_course_id = c.zyk_id
    WHERE e.zyk_student_id = NEW.zyk_student_id AND e.zyk_score IS NOT NULL;

    IF zyk_total_credits > 0 THEN
        zyk_calculated_gpa := zyk_total_grade_points / zyk_total_credits;
    ELSE
        zyk_calculated_gpa := 0;
    END IF;

    UPDATE Zhaoyk_stu
    SET zyk_gpa = zyk_calculated_gpa,
        zyk_credits = zyk_total_credits,
        zyk_update_time = CURRENT_TIMESTAMP
    WHERE zyk_student_id = NEW.zyk_student_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- score 表触发器
DROP TRIGGER IF EXISTS Zhaoyk_trg_update_gpa_after_score ON Zhaoyk_score;
CREATE TRIGGER Zhaoyk_trg_update_gpa_after_score
AFTER INSERT OR UPDATE OF zyk_score ON Zhaoyk_score
FOR EACH ROW
EXECUTE PROCEDURE Zhaoyk_update_student_gpa();

-- 2.3 score 表时间检测触发器
CREATE OR REPLACE FUNCTION Zhaoyk_check_score_time_validity()
RETURNS TRIGGER AS $$
DECLARE
    zyk_course_created_at TIMESTAMP;
BEGIN
    SELECT zyk_create_time INTO zyk_course_created_at
    FROM Zhaoyk_cou
    WHERE zyk_id = NEW.zyk_course_id;

    IF NEW.zyk_create_time < zyk_course_created_at THEN
        RAISE EXCEPTION  '成绩记录创建时间不能早于课程创建时间: % < %', NEW.zyk_create_time, zyk_course_created_at;
    END IF;

    IF NEW.zyk_update_time < zyk_course_created_at THEN
        RAISE EXCEPTION '成绩记录更新时间不能早于课程创建时间: % < %', NEW.zyk_update_time, zyk_course_created_at;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- score 表触发器
DROP TRIGGER IF EXISTS Zhaoyk_trg_check_score_time ON Zhaoyk_score;
CREATE TRIGGER Zhaoyk_trg_check_score_time
BEFORE INSERT OR UPDATE OF zyk_create_time, zyk_update_time ON Zhaoyk_score
FOR EACH ROW
EXECUTE PROCEDURE Zhaoyk_check_score_time_validity();

-- 2.4 学生密码生成触发器
CREATE OR REPLACE FUNCTION Zhaoyk_set_default_student_password()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.zyk_password IS NULL OR TRIM(NEW.zyk_password) = '' THEN
        NEW.zyk_password := 'zjut' || RIGHT(NEW.zyk_student_id, 6);
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- student 表触发器
DROP TRIGGER IF EXISTS Zhaoyk_trg_default_password_student ON Zhaoyk_stu;
CREATE TRIGGER Zhaoyk_trg_default_password_student
BEFORE INSERT ON Zhaoyk_stu
FOR EACH ROW
EXECUTE PROCEDURE Zhaoyk_set_default_student_password();

-- 3. 视图设计
-- 3.1 查询"每门课程每学年的平均成绩"
DROP VIEW IF EXISTS Zhaoyk_course_avg_score_view;
CREATE VIEW Zhaoyk_course_avg_score_view AS
SELECT
    Zhaoyk_cou.zyk_id AS zyk_course_id,
    Zhaoyk_cou.zyk_name AS zyk_course_name,
    Zhaoyk_cou.zyk_school_year,
    ROUND(AVG(Zhaoyk_score.zyk_score), 2) AS zyk_avg_score
FROM Zhaoyk_cou
JOIN Zhaoyk_score ON Zhaoyk_cou.zyk_id = Zhaoyk_score.zyk_course_id
WHERE Zhaoyk_score.zyk_score IS NOT NULL
GROUP BY Zhaoyk_cou.zyk_id, Zhaoyk_cou.zyk_name, Zhaoyk_cou.zyk_school_year
ORDER BY Zhaoyk_cou.zyk_name, Zhaoyk_cou.zyk_school_year;

-- 3.2 查询"各生源地招收学生数量"
DROP VIEW IF EXISTS Zhaoyk_province_city_student_count_view;
CREATE VIEW Zhaoyk_province_city_student_count_view AS
SELECT
    Zhaoyk_province.zyk_id AS zyk_province_id,
    Zhaoyk_province.zyk_name AS zyk_province_name,
    Zhaoyk_city.zyk_id AS zyk_city_id,
    Zhaoyk_city.zyk_name AS zyk_city_name,
    COUNT(Zhaoyk_stu.zyk_id) AS zyk_student_count
FROM Zhaoyk_stu
JOIN Zhaoyk_city ON Zhaoyk_stu.zyk_city_id = Zhaoyk_city.zyk_id
JOIN Zhaoyk_province ON Zhaoyk_city.zyk_province_id = Zhaoyk_province.zyk_id
GROUP BY Zhaoyk_province.zyk_id, Zhaoyk_province.zyk_name, Zhaoyk_city.zyk_id, Zhaoyk_city.zyk_name
ORDER BY Zhaoyk_province.zyk_name, Zhaoyk_city.zyk_name;

-- 3.3 查询"各生源地在各专业绩点前十的数量"
DROP VIEW IF EXISTS Zhaoyk_top10_gpa_by_province_department_view;
CREATE VIEW Zhaoyk_top10_gpa_by_province_department_view AS
WITH ranked_students AS (
    SELECT
        Zhaoyk_stu.zyk_id,
        Zhaoyk_stu.zyk_name AS zyk_student_name,
        Zhaoyk_department.zyk_name AS zyk_department_name,
        Zhaoyk_province.zyk_name AS zyk_province_name,
        Zhaoyk_stu.zyk_gpa,
        RANK() OVER (PARTITION BY Zhaoyk_department.zyk_id, Zhaoyk_province.zyk_id ORDER BY Zhaoyk_stu.zyk_gpa DESC) AS zyk_rank
    FROM Zhaoyk_stu
    JOIN Zhaoyk_class ON Zhaoyk_stu.zyk_class_id = Zhaoyk_class.zyk_id
    JOIN Zhaoyk_department ON Zhaoyk_class.zyk_department_id = Zhaoyk_department.zyk_id
    JOIN Zhaoyk_city ON Zhaoyk_stu.zyk_city_id = Zhaoyk_city.zyk_id
    JOIN Zhaoyk_province ON Zhaoyk_city.zyk_province_id = Zhaoyk_province.zyk_id
)
SELECT
    zyk_province_name,
    zyk_department_name,
    COUNT(*) AS zyk_top10_count
FROM ranked_students
WHERE zyk_rank <= 10
GROUP BY zyk_province_name, zyk_department_name
ORDER BY zyk_province_name, zyk_department_name;

-- 3.4 查询"每位教师的每门课程的学生人数和平均成绩"
DROP VIEW IF EXISTS Zhaoyk_teacher_course_view;
CREATE VIEW Zhaoyk_teacher_course_view AS
SELECT
    Zhaoyk_tea.zyk_teacher_id,
    Zhaoyk_tea.zyk_name AS zyk_teacher_name,
    Zhaoyk_cou.zyk_id AS zyk_course_id,
    Zhaoyk_cou.zyk_name AS zyk_course_name,
    Zhaoyk_cou.zyk_school_year,
    Zhaoyk_cou.zyk_semester,
    COUNT(DISTINCT Zhaoyk_score.zyk_student_id) AS zyk_student_count,
    ROUND(AVG(Zhaoyk_score.zyk_score), 2) AS zyk_avg_score
FROM Zhaoyk_tea
JOIN Zhaoyk_teaching ON Zhaoyk_tea.zyk_teacher_id = Zhaoyk_teaching.zyk_teacher_id
JOIN Zhaoyk_cou ON Zhaoyk_teaching.zyk_course_id = Zhaoyk_cou.zyk_id
LEFT JOIN Zhaoyk_score ON Zhaoyk_score.zyk_course_id = Zhaoyk_cou.zyk_id
GROUP BY Zhaoyk_tea.zyk_teacher_id, Zhaoyk_tea.zyk_name, Zhaoyk_cou.zyk_id, Zhaoyk_cou.zyk_name, Zhaoyk_cou.zyk_school_year, Zhaoyk_cou.zyk_semester
ORDER BY Zhaoyk_tea.zyk_teacher_id, Zhaoyk_cou.zyk_school_year, Zhaoyk_cou.zyk_name;

-- 4. 存储过程
-- 4.1 课程成绩分布（按成绩分段统计某门课某学年的成绩分布）
CREATE OR REPLACE FUNCTION Zhaoyk_select_course_score_distribution(
    zyk_cid INT,
    zyk_schoolyear INT
)
    RETURNS TABLE(zyk_score_range TEXT, zyk_count BIGINT)
    STRICT
AS
$$
BEGIN
    RETURN QUERY
        SELECT
            (width_bucket(Zhaoyk_score.zyk_score, 0, 100, 10)-1)*10 || '-' || (width_bucket(Zhaoyk_score.zyk_score, 0, 100, 10)*10 - 1) AS zyk_score_range,
            COUNT(*) AS zyk_count
        FROM
            Zhaoyk_score
                JOIN Zhaoyk_cou ON Zhaoyk_score.zyk_course_id = Zhaoyk_cou.zyk_id
        WHERE
            Zhaoyk_cou.zyk_id = zyk_cid
          AND Zhaoyk_cou.zyk_school_year = zyk_schoolyear
          AND Zhaoyk_score.zyk_score IS NOT NULL
        GROUP BY
            zyk_score_range
        ORDER BY
            zyk_score_range;
END;
$$ LANGUAGE plpgsql;

-- 4.2 按专业学生总绩点排名
CREATE OR REPLACE FUNCTION Zhaoyk_select_department_gpa_ranking(
    zyk_dept_id INT
)
RETURNS TABLE(
    zyk_rank BIGINT,
    zyk_student_id CHAR(12),
    zyk_student_name VARCHAR,
    zyk_class_id INT,
    zyk_class_name VARCHAR,
    zyk_gpa NUMERIC(4,2)
)
AS
$$
BEGIN
    RETURN QUERY
    SELECT
        RANK() OVER (ORDER BY Zhaoyk_stu.zyk_gpa DESC) AS zyk_rank,
        Zhaoyk_stu.zyk_student_id,
        Zhaoyk_stu.zyk_name AS zyk_student_name,
        Zhaoyk_class.zyk_id AS zyk_class_id,
        Zhaoyk_class.zyk_name AS zyk_class_name,
        Zhaoyk_stu.zyk_gpa
    FROM
        Zhaoyk_stu
    JOIN Zhaoyk_class ON Zhaoyk_stu.zyk_class_id = Zhaoyk_class.zyk_id
    WHERE
        Zhaoyk_class.zyk_department_id = zyk_dept_id;
END;
$$ LANGUAGE plpgsql;

-- 4.3 按班级学生总绩点排名
CREATE OR REPLACE FUNCTION Zhaoyk_select_class_gpa_ranking(
    zyk_classId INT
)
    RETURNS TABLE (
                      zyk_rank BIGINT,
                      zyk_student_id CHAR(12),
                      zyk_student_name VARCHAR,
                      zyk_class_id INT,
                      zyk_class_name VARCHAR,
                      zyk_gpa NUMERIC(4,2)
                  )
AS
$$
BEGIN
    RETURN QUERY
        SELECT
                    RANK() OVER (ORDER BY Zhaoyk_stu.zyk_gpa DESC) AS zyk_rank,
                    Zhaoyk_stu.zyk_student_id,
                    Zhaoyk_stu.zyk_name,
                    Zhaoyk_class.zyk_id AS zyk_class_id,
                    Zhaoyk_class.zyk_name AS zyk_class_name,
                    Zhaoyk_stu.zyk_gpa
        FROM
            Zhaoyk_stu
                JOIN Zhaoyk_class ON Zhaoyk_stu.zyk_class_id = Zhaoyk_class.zyk_id
        WHERE
            Zhaoyk_stu.zyk_class_id = zyk_classId
          AND Zhaoyk_stu.zyk_gpa IS NOT NULL;
END;
$$ LANGUAGE plpgsql;

-- 4.4 查询学生所有课程与成绩及是否需要重修
CREATE OR REPLACE FUNCTION Zhaoyk_select_student_scores_and_retake_status(
    zyk_sid VARCHAR,
    zyk_schoolyear INT
)
    RETURNS TABLE(
                     zyk_course_id INT,
                     zyk_course_name VARCHAR,
                     zyk_credit NUMERIC(4,2),
                     zyk_score NUMERIC(5,2),
                     zyk_retake_required VARCHAR
                 ) AS
$$
BEGIN
    RETURN QUERY
        SELECT
            Zhaoyk_cou.zyk_id,
            Zhaoyk_cou.zyk_name,
            Zhaoyk_cou.zyk_credit,
            Zhaoyk_score.zyk_score,
            CASE
                WHEN abs(extract(epoch from Zhaoyk_score.zyk_update_time - Zhaoyk_score.zyk_create_time)) < 1 THEN '-'::VARCHAR
                WHEN Zhaoyk_score.zyk_score < 60 THEN '是'::VARCHAR
                ELSE '否'::VARCHAR
                END AS zyk_retake_required
        FROM
            Zhaoyk_score
                JOIN Zhaoyk_cou ON Zhaoyk_score.zyk_course_id = Zhaoyk_cou.zyk_id
        WHERE
            Zhaoyk_score.zyk_student_id = zyk_sid
          AND Zhaoyk_cou.zyk_school_year = zyk_schoolyear;
END;
$$ LANGUAGE plpgsql;

-- 4.5 查询教师教授课程
CREATE OR REPLACE FUNCTION Zhaoyk_select_teacher_courses(
    zyk_tid INT,
    zyk_in_school_year INT,
    zyk_in_semester VARCHAR(1)
)
    RETURNS TABLE(
                     zyk_teacher_id VARCHAR(10),
                     zyk_teacher_name VARCHAR,
                     zyk_course_id INT,
                     zyk_course_name VARCHAR,
                     zyk_class_id INT,
                     zyk_class_name VARCHAR,
                     zyk_school_year INT,
                     zyk_credit NUMERIC(4,2),
                     zyk_hours INT,
                     zyk_exam_type VARCHAR,
                     zyk_semester VARCHAR(1)
                     
                 ) AS
$$
BEGIN
    RETURN QUERY
        SELECT
            Zhaoyk_tea.zyk_teacher_id,
            Zhaoyk_tea.zyk_name,
            Zhaoyk_cou.zyk_id,
            Zhaoyk_cou.zyk_name,
            Zhaoyk_class.zyk_id,
            Zhaoyk_class.zyk_name,
            Zhaoyk_cou.zyk_school_year,
            Zhaoyk_cou.zyk_credit,
            Zhaoyk_cou.zyk_hours,
            Zhaoyk_cou.zyk_exam_type,
            Zhaoyk_cou.zyk_semester::VARCHAR
        FROM
            Zhaoyk_teaching
                JOIN Zhaoyk_tea ON Zhaoyk_teaching.zyk_teacher_id = Zhaoyk_tea.zyk_teacher_id
                JOIN Zhaoyk_cou ON Zhaoyk_teaching.zyk_course_id = Zhaoyk_cou.zyk_id
                JOIN Zhaoyk_class ON Zhaoyk_cou.zyk_class_id = Zhaoyk_class.zyk_id
        WHERE
            Zhaoyk_tea.zyk_id = zyk_tid
          AND Zhaoyk_cou.zyk_school_year = zyk_in_school_year
          AND Zhaoyk_cou.zyk_semester = zyk_in_semester;
END;
$$ LANGUAGE plpgsql;

-- 4.6 查询教师各课程的平均成绩
CREATE OR REPLACE FUNCTION Zhaoyk_select_teacher_course_avg_scores(
    zyk_tid INT
)
RETURNS TABLE(
    zyk_teacher_id VARCHAR(10),
    zyk_course_name VARCHAR,
    zyk_class_name VARCHAR,
    zyk_school_year INT,
    zyk_avg_score NUMERIC(5,2)
) AS
$$
BEGIN
    RETURN QUERY
    SELECT
        Zhaoyk_tea.zyk_teacher_id,
        Zhaoyk_cou.zyk_name,
        Zhaoyk_class.zyk_name,
        Zhaoyk_cou.zyk_school_year,
        ROUND(AVG(Zhaoyk_score.zyk_score), 2)
    FROM
        Zhaoyk_tea
    JOIN Zhaoyk_teaching ON Zhaoyk_tea.zyk_teacher_id = Zhaoyk_teaching.zyk_teacher_id
    JOIN Zhaoyk_cou ON Zhaoyk_teaching.zyk_course_id = Zhaoyk_cou.zyk_id
    JOIN Zhaoyk_score ON Zhaoyk_cou.zyk_id = Zhaoyk_score.zyk_course_id
    JOIN Zhaoyk_class ON Zhaoyk_cou.zyk_class_id = Zhaoyk_class.zyk_id
    WHERE Zhaoyk_tea.zyk_id = zyk_tid
    GROUP BY Zhaoyk_tea.zyk_teacher_id, Zhaoyk_cou.zyk_name, Zhaoyk_class.zyk_name, Zhaoyk_cou.zyk_school_year;
END;
$$ LANGUAGE plpgsql;

-- 4.7 查询某学年某个相同课程名的课程的所有学生成绩总览
CREATE OR REPLACE FUNCTION Zhaoyk_select_course_scores_by_name_and_year(zyk_coursename VARCHAR, zyk_schoolyear INT)
    RETURNS TABLE(zyk_student_id INT, zyk_student_name VARCHAR, zyk_course_name VARCHAR, zyk_class_name VARCHAR, zyk_semester CHAR, zyk_teacher_name VARCHAR, zyk_score NUMERIC)
    LANGUAGE plpgsql
AS
$$
BEGIN
    RETURN QUERY
        SELECT
            Zhaoyk_stu.zyk_id,
            Zhaoyk_stu.zyk_name,
            Zhaoyk_cou.zyk_name,
            Zhaoyk_class.zyk_name,
            Zhaoyk_cou.zyk_semester,
            Zhaoyk_tea.zyk_name,
            Zhaoyk_score.zyk_score
        FROM
            Zhaoyk_cou
                JOIN Zhaoyk_teaching ON Zhaoyk_cou.zyk_id = Zhaoyk_teaching.zyk_course_id
                JOIN Zhaoyk_tea ON Zhaoyk_teaching.zyk_teacher_id = Zhaoyk_tea.zyk_teacher_id
                JOIN Zhaoyk_score ON Zhaoyk_score.zyk_course_id = Zhaoyk_cou.zyk_id
                JOIN Zhaoyk_stu ON Zhaoyk_score.zyk_student_id = Zhaoyk_stu.zyk_student_id
                JOIN Zhaoyk_class ON Zhaoyk_stu.zyk_class_id = Zhaoyk_class.zyk_id
        WHERE
            Zhaoyk_cou.zyk_school_year = zyk_schoolyear
          AND Zhaoyk_cou.zyk_name ILIKE '%' || zyk_coursename || '%';
END;
$$;

-- 4.8 查询某个班级的课程表
CREATE OR REPLACE FUNCTION Zhaoyk_select_class_course_schedule(
    zyk_in_class_id INT,
    zyk_in_school_year INT,
    zyk_in_semester VARCHAR(1)
)
    RETURNS TABLE (
                      zyk_course_id INT,
                      zyk_course_name VARCHAR,
                      zyk_school_year INT,
                      zyk_semester CHAR(1),
                      zyk_credit NUMERIC(4,2),
                      zyk_hours INT,
                      zyk_exam_type VARCHAR,
                      zyk_teacher_name VARCHAR,
                      zyk_assign_time TIMESTAMP
                  ) AS
$$
BEGIN
    RETURN QUERY
        SELECT
            Zhaoyk_cou.zyk_id AS zyk_course_id,
            Zhaoyk_cou.zyk_name AS zyk_course_name,
            Zhaoyk_cou.zyk_school_year,
            Zhaoyk_cou.zyk_semester,
            Zhaoyk_cou.zyk_credit,
            Zhaoyk_cou.zyk_hours,
            Zhaoyk_cou.zyk_exam_type,
            Zhaoyk_tea.zyk_name AS zyk_teacher_name,
            Zhaoyk_teaching.zyk_assign_time
        FROM
            Zhaoyk_cou
                JOIN Zhaoyk_teaching ON Zhaoyk_cou.zyk_id = Zhaoyk_teaching.zyk_course_id
                JOIN Zhaoyk_tea ON Zhaoyk_teaching.zyk_teacher_id = Zhaoyk_tea.zyk_teacher_id
        WHERE
            Zhaoyk_cou.zyk_class_id = zyk_in_class_id
          AND Zhaoyk_cou.zyk_school_year = zyk_in_school_year
          AND Zhaoyk_cou.zyk_semester = zyk_in_semester;
END;
$$ LANGUAGE plpgsql;

-- 4.9 查询某位老师所教授的某门课程下的所有学生列表
CREATE OR REPLACE FUNCTION Zhaoyk_select_students_by_teacher_course(
    zyk_in_teacher_id VARCHAR(10),
    zyk_in_course_id INT
)
    RETURNS TABLE (
                      zyk_student_id INT,
                      zyk_student_name VARCHAR,
                      zyk_class_name VARCHAR,
                      zyk_score NUMERIC(5,2),
                      zyk_rank BIGINT
                  ) AS
$$
BEGIN
    RETURN QUERY
        SELECT
            Zhaoyk_stu.zyk_id,
            Zhaoyk_stu.zyk_name,
            Zhaoyk_class.zyk_name AS zyk_class_name,
            Zhaoyk_score.zyk_score,
            RANK() OVER (ORDER BY Zhaoyk_score.zyk_score DESC NULLS LAST) AS zyk_rank
        FROM
            Zhaoyk_teaching
                JOIN Zhaoyk_cou ON Zhaoyk_teaching.zyk_course_id = Zhaoyk_cou.zyk_id
                JOIN Zhaoyk_score ON Zhaoyk_cou.zyk_id = Zhaoyk_score.zyk_course_id
                JOIN Zhaoyk_stu ON Zhaoyk_stu.zyk_student_id = Zhaoyk_score.zyk_student_id
                JOIN Zhaoyk_class ON Zhaoyk_stu.zyk_class_id = Zhaoyk_class.zyk_id
        WHERE
            Zhaoyk_teaching.zyk_teacher_id = zyk_in_teacher_id
          AND Zhaoyk_cou.zyk_id = zyk_in_course_id;
END;
$$ LANGUAGE plpgsql;

-- 5. 数据准备
-- 1. 插入 department 表
INSERT INTO Zhaoyk_department (zyk_name) VALUES
('计算机科学与技术'),
('电子信息工程'),
('机械设计制造及其自动化'),
('土木工程'),
('金融学'),
('汉语言文学'),
('法学'),
('环境工程'),
('数学与应用数学'),
('英语');

-- 2. 插入 class 表
INSERT INTO Zhaoyk_class (zyk_name, zyk_department_id) VALUES
('计算机2303班', 1),
('计算机2304班', 1),
('电子信息2303班', 2),
('电子信息2304班', 2),
('金融2303班', 5),
('金融2304班', 5),
('汉语言2303班', 6),
('汉语言2304班', 6),
('法学2303班', 7),
('法学2304班', 7);

-- 3. 插入 student 表
INSERT INTO Zhaoyk_stu (zyk_student_id, zyk_password, zyk_name, zyk_gender, zyk_age, zyk_city_id, zyk_class_id, zyk_gpa) VALUES
('302023315101', 'zjut315101', '张三', 'M', 20, 1, 1, 0),
('302023315102', 'zjut315102', '李四', 'F', 21, 2, 2, 0),
('302023315103', 'zjut315103', '王五', 'M', 19, 3, 3, 0),
('302023315104', 'zjut315104', '赵六', 'F', 22, 4, 4, 0),
('302023315105', 'zjut315105', '钱七', 'M', 23, 5, 5, 0),
('302023315106', 'zjut315106', '孙八', 'F', 24, 6, 6, 0),
('302023315107', 'zjut315107', '周九', 'M', 20, 7, 7, 0),
('302023315108', 'zjut315108', '吴十', 'F', 21, 8, 8, 0),
('302023315109', 'zjut315109', '郑十一', 'M', 22, 9, 9, 0),
('302023315110', 'zjut315110', '冯十二', 'F', 23, 10, 10, 0),
('302023315111', 'zjut315111', '王丹', 'M', 18, 5, 4, 0),
('302023315112', 'zjut315112', '刘旭', 'M', 23, 9, 2, 0),
('302023315113', 'zjut315113', '李静', 'F', 22, 6, 7, 0),
('302023315114', 'zjut315114', '陈伟', 'M', 21, 3, 1, 0),
('302023315115', 'zjut315115', '杨磊', 'M', 20, 1, 5, 0),
('302023315116', 'zjut315116', '徐婷', 'F', 19, 10, 6, 0),
('302023315117', 'zjut315117', '周娜', 'F', 20, 8, 9, 0),
('302023315118', 'zjut315118', '吴峰', 'M', 22, 2, 8, 0),
('302023315119', 'zjut315119', '张洋', 'M', 24, 4, 3, 0),
('302023315120', 'zjut315120', '林倩', 'F', 21, 7, 10, 0);

-- 4. 插入 teacher 表
INSERT INTO Zhaoyk_tea (zyk_teacher_id, zyk_password, zyk_name, zyk_gender, zyk_age, zyk_title, zyk_phone, zyk_is_admin) VALUES
('10000001', 'zjut123', '王老师', 'M', 45, '副教授', '13800000001', TRUE),
('10000002', 'zjut123', '李老师', 'F', 38, '讲师', '13800000002', FALSE),
('10000003', 'zjut123', '赵老师', 'M', 50, '教授', '13800000003', FALSE),
('10000004', 'zjut123', '钱老师', 'F', 43, '副教授', '13800000004', FALSE),
('10000005', 'zjut123', '孙老师', 'M', 39, '讲师', '13800000005', FALSE),
('10000006', 'zjut123', '周老师', 'F', 55, '教授', '13800000006', FALSE),
('10000007', 'zjut123', '吴老师', 'M', 42, '讲师', '13800000007', FALSE),
('10000008', 'zjut123', '郑老师', 'F', 37, '副教授', '13800000008', FALSE),
('10000009', 'zjut123', '冯老师', 'M', 60, '教授', '13800000009', FALSE),
('10000010', 'zjut123', '陈老师', 'F', 46, '讲师', '13800000010', FALSE),
('10000011', 'zjut123', '毛燕', 'F', 48, '讲师', '13800000011', FALSE),
('10000012', 'zjut123', '梁波', 'M', 31, '教授', '13800000012', FALSE),
('10000013', 'zjut123', '刘洁', 'F', 37, '副教授', '13800000013', FALSE),
('10000014', 'zjut123', '赵雷', 'M', 55, '讲师', '13800000014', FALSE),
('10000015', 'zjut123', '唐慧', 'F', 41, '副教授', '13800000015', FALSE),
('10000016', 'zjut123', '郭强', 'M', 50, '教授', '13800000016', FALSE),
('10000017', 'zjut123', '沈丽', 'F', 36, '讲师', '13800000017', FALSE),
('10000018', 'zjut123', '曹洋', 'M', 40, '副教授', '13800000018', FALSE),
('10000019', 'zjut123', '邓婕', 'F', 35, '教授', '13800000019', FALSE),
('10000020', 'zjut123', '傅林', 'M', 60, '讲师', '13800000020', FALSE);

-- 5. 插入 course 表
INSERT INTO Zhaoyk_cou (zyk_name, zyk_semester, zyk_hours, zyk_credit, zyk_class_id, zyk_exam_type, zyk_school_year)
VALUES ('数据库原理', '1', 48, 3.0, 1, '考试', 2025),
       ('操作系统', '2', 64, 3.5, 2, '考试', 2025),
       ('电路基础', '1', 56, 3.0, 3, '考查', 2025),
       ('混凝土结构', '2', 40, 2.5, 4, '考试', 2025),
       ('金融市场', '1', 36, 2.0, 5, '考查', 2025),
       ('古代文学', '2', 32, 2.5, 6, '考试', 2025),
       ('宪法学', '1', 48, 3.0, 7, '考查', 2025),
       ('水质工程', '2', 60, 3.0, 8, '考试', 2025),
       ('数学分析', '1', 64, 4.0, 9, '考试', 2025),
       ('英语听说', '2', 36, 2.0, 10, '考查', 2025),
       ('保险核保', '1', 42, 2.4, 5, '考试', 2025),
       ('英语翻译', '2', 35, 1.7, 1, '考查', 2025),
       ('数据结构', '1', 48, 3.0, 1, '考试', 2025),
       ('电磁场理论', '2', 56, 3.5, 3, '考试', 2025),
       ('建筑工程识图', '1', 36, 2.0, 4, '考查', 2025),
       ('财务管理', '2', 40, 2.5, 5, '考试', 2025),
       ('现代汉语', '1', 48, 3.0, 6, '考试', 2025),
       ('民法总论', '2', 60, 3.0, 7, '考试', 2025),
       ('环境科学概论', '1', 40, 2.0, 8, '考查', 2025),
       ('数学建模', '2', 60, 4.0, 9, '考试', 2025);

-- 6. 插入 teaching 表
INSERT INTO Zhaoyk_teaching (zyk_teacher_id, zyk_course_id)
VALUES ('10000001', 1),
       ('10000002', 2),
       ('10000003', 3),
       ('10000004', 4),
       ('10000005', 5),
       ('10000006', 6),
       ('10000007', 7),
       ('10000008', 8),
       ('10000009', 9),
       ('10000010', 10),
       ('10000018', 5),
       ('10000016', 18),
       ('10000015', 14),
       ('10000014', 12),
       ('10000013', 17),
       ('10000012', 13),
       ('10000011', 16),
       ('10000020', 11),
       ('10000019', 19),
       ('10000017', 20);

-- 7. 插入 province 表
INSERT INTO Zhaoyk_province (zyk_name) VALUES
('浙江省'), ('江苏省'), ('广东省'), ('山东省'), ('河南省'),
('四川省'), ('湖北省'), ('湖南省'), ('北京市'), ('上海市');

-- 8. 插入 city 表
INSERT INTO Zhaoyk_city (zyk_name, zyk_province_id) VALUES
('杭州', 1), ('南京', 2), ('广州', 3), ('济南', 4), ('郑州', 5),
('成都', 6), ('武汉', 7), ('长沙', 8), ('北京', 9), ('上海', 10);

-- 9. 插入 score 表
INSERT INTO Zhaoyk_score (zyk_student_id, zyk_course_id, zyk_score)
VALUES 
('302023315101', 1, NULL),
('302023315101', 13, 95.00),
('302023315101', 12, NULL),
('302023315114', 1, NULL),
('302023315114', 13, NULL),
('302023315114', 12, NULL),
('302023315102', 2, NULL),
('302023315112', 2, NULL),
('302023315103', 3, NULL),
('302023315103', 14, NULL),
('302023315119', 3, NULL),
('302023315119', 14, NULL),
('302023315104', 4, NULL),
('302023315104', 15, NULL),
('302023315111', 4, NULL),
('302023315111', 15, NULL),
('302023315105', 5, NULL),
('302023315105', 11, NULL),
('302023315105', 16, NULL),
('302023315115', 5, NULL),
('302023315115', 11, NULL),
('302023315115', 16, NULL),
('302023315106', 6, NULL),
('302023315116', 6, NULL),
('302023315107', 7, NULL),
('302023315107', 17, NULL),
('302023315113', 7, NULL),
('302023315113', 17, NULL),
('302023315108', 8, NULL),
('302023315118', 8, NULL),
('302023315109', 9, NULL),
('302023315109', 18, NULL),
('302023315117', 9, NULL),
('302023315117', 18, NULL),
('302023315110', 10, NULL),
('302023315110', 19, NULL),
('302023315110', 20, NULL),
('302023315120', 10, NULL),
('302023315120', 19, NULL),
('302023315120', 20, NULL);