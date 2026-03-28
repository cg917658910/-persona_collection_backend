ALTER TABLE public.pm_themes
ADD COLUMN IF NOT EXISTS subject_type text NOT NULL DEFAULT 'character';

UPDATE public.pm_themes
SET subject_type = 'character'
WHERE COALESCE(btrim(subject_type), '') = '';

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'chk_pm_themes_subject_type'
  ) THEN
    ALTER TABLE public.pm_themes
    ADD CONSTRAINT chk_pm_themes_subject_type
    CHECK (subject_type IN ('character', 'relation'));
  END IF;
END $$;
