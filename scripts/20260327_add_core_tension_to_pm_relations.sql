ALTER TABLE public.pm_relations
ADD COLUMN IF NOT EXISTS core_tension varchar(255);

COMMENT ON COLUMN public.pm_relations.core_tension IS '一条短句概括这段关系最深的张力';
