drop function update_job_search_vector();
;;
create function update_job_search_vector() returns trigger as $$
begin
    new.search_vector :=
      setweight(to_tsvector('pg_catalog.norwegian', coalesce(NEW.description_value, '')), 'A') ||
      setweight(to_tsvector('pg_catalog.norwegian', coalesce(array_to_string(NEW.tags, ' '), '')), 'A') ||
      setweight(to_tsvector('pg_catalog.norwegian', coalesce(NEW.address_city, '')), 'B') ||
      setweight(to_tsvector('pg_catalog.norwegian', coalesce(NEW.address_district, '')), 'B') ||
      setweight(to_tsvector('pg_catalog.norwegian', coalesce(NEW.address_postal_code, '')), 'B') ||
      setweight(to_tsvector('pg_catalog.norwegian', coalesce(NEW.address_line_1, '')), 'C') ||
      setweight(to_tsvector('pg_catalog.norwegian', coalesce(NEW.address_line_2, '')), 'C');
return new;
end
$$ language plpgsql;
;;
create trigger tsvectorupdate before insert or update on job for each row execute procedure update_job_search_vector();;
