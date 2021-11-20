-- jobs.jobList definition

CREATE TABLE `jobList` (
  `job_id` varchar(30) NOT NULL,
  `name` varchar(100) DEFAULT NULL,
  `created_at` datetime NOT NULL,
  `created_by` varchar(100) DEFAULT NULL,
  `modified_at` datetime NOT NULL,
  `modified_by` varchar(100) DEFAULT NULL,
  `src_url` varchar(1000) DEFAULT NULL,
  `status` varchar(15) NOT NULL,
  `error_msg` varchar(2000) DEFAULT NULL,
  `tech_info` mediumtext
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
