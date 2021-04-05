FROM webdevops/liquibase:postgres

COPY . .
RUN printf 'parameter.schema: public\n parameter.tablespace_name: pg_default \n' >> liquibase.properties
