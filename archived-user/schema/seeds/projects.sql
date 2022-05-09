INSERT INTO projects (project_id, name, prefix, description, user_id, team_id, active, public, column_order, updated_at, created_at) VALUES
('8695a94f-7e0a-4198-8c0a-d3e12727a5ba','Testing','Tes-','A test','a4b54ec1-57f9-4c39-ab53-d936dbb6c177','',true,false,'{"column-1","column-2","column-3","column-4"}','2021-05-16T21:02:51.371324Z','2021-05-16T07:36:09.654601Z')
ON CONFLICT DO NOTHING;
