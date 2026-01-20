create or replace function update_job_search_vector() returns trigger as $$
begin
    new.search_vector :=
      setweight(to_tsvector('pg_catalog.norwegian', coalesce(new.description_value, '')), 'a') ||
      setweight(to_tsvector('pg_catalog.norwegian', coalesce(array_to_string(new.tags, ' '), '')), 'a') ||
      setweight(to_tsvector('pg_catalog.norwegian', coalesce(new.address_city, '')), 'b') ||
      setweight(to_tsvector('pg_catalog.norwegian', coalesce(new.address_district, '')), 'b') ||
      setweight(to_tsvector('pg_catalog.norwegian', coalesce(new.address_postal_code, '')), 'b') ||
      setweight(to_tsvector('pg_catalog.norwegian', coalesce(new.address_line_1, '')), 'c') ||
      setweight(to_tsvector('pg_catalog.norwegian', coalesce(new.address_line_2, '')), 'c');
return new;
end
$$ language plpgsql
;;
create trigger job_trigger_search_vector before insert or update on job for each row execute procedure update_job_search_vector()
;;
create or replace function update_job_earth() returns trigger as $$
begin
  new.earth_point := ll_to_earth(new.position_latitude, new.position_longitude);
return new;
end;
$$ language plpgsql
;;
CREATE TRIGGER job_trigger_earth BEFORE INSERT OR UPDATE OF position_latitude, position_longitude ON job FOR EACH ROW EXECUTE FUNCTION update_job_earth();


