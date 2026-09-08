-- ============================================================================
-- seed.sql  —  "Kusuriya no Hitorigoto" (The Apothecary Diaries) diagnosis toy
-- ----------------------------------------------------------------------------
-- A Bayesian-style seed: causes (poisons/diseases) carry a base `prior`,
-- symptom_profiles hold P(symptom | cause) as `likelihood`, and
-- factor_profiles hold P(context-factor | cause) as `modifier` (<= 1).
--
-- No IDs are hardcoded (columns are GENERATED ALWAYS AS IDENTITY). Every FK is
-- resolved via the UNIQUE natural keys: remedies.name, causes.code,
-- factors.code, symptoms.code, categories.name (unique within this seed),
-- cases.title. Safe to re-run after TRUNCATE ... RESTART IDENTITY CASCADE.
-- ============================================================================

BEGIN;

-- ----------------------------------------------------------------------------
-- Remedies (Maomao's cures and antidotes)
-- ----------------------------------------------------------------------------
INSERT INTO remedies (name, description, notes) VALUES
  ('Rice Bran Regimen',        'Unpolished rice and bran to restore the exhausted body.', 'Maomao''s answer to the wasting sickness of the affluent.'),
  ('Activated Charcoal',       'Adsorbs ingested poison before it spreads through the body.', 'Most useful shortly after the poison is swallowed.'),
  ('Purge and Cease Powder',   'Discontinue the white face powder and purge the metal from the body.', 'The cosmetic itself is the poison.'),
  ('Cinchona Bark',            'Bitter bark decoction to break the recurring fever.', 'A foreign remedy for the shaking fever.'),
  ('Ginger Infusion',          'Warm ginger to settle a turning stomach.', NULL),
  ('Antivenom Serum',          'Serum drawn against the offending venom.', 'Only as good as the snake it was matched to.'),
  ('Citrus and Greens',        'Fresh fruit and vegetables to heal bleeding gums.', 'The sailor''s and the shut-in''s cure.'),
  ('Rehydration Broth',        'Salted water and broth to replace what the body has lost.', 'Fluid loss kills faster than the illness itself.'),
  ('Isolation and Rest',       'Quarantine and supportive care until the fever passes.', 'Protects the household as much as the patient.'),
  ('Emetic Draught',           'Induce vomiting to expel a freshly swallowed poison.', 'Useless once the poison has taken hold.'),
  ('Palace Physician''s Antidote', 'A prepared counter-agent for slowed heart and clouded mind.', 'Kept close at every banquet.'),
  ('Mercury Purge',            'A harsh mineral course; the treatment borders on the poison.', 'Historically dosed for the courtesans'' pox.'),
  ('Liver and Iron Tonic',     'Iron-rich fare to rebuild thin, pale blood.', NULL),
  ('Stimulant Draught',        'Rouses breathing and pulse suppressed by the drug.', 'For the poppy''s deep sleep.');

-- ----------------------------------------------------------------------------
-- Symptom categories
-- ----------------------------------------------------------------------------
INSERT INTO categories (name) VALUES
  ('Neurological'),
  ('Gastrointestinal'),
  ('Dermatological'),
  ('Cardiovascular'),
  ('Respiratory'),
  ('Constitutional'),
  ('Ophthalmic'),
  ('Psychiatric');

-- ----------------------------------------------------------------------------
-- Symptoms (description kept <= 50 chars per schema)
-- ----------------------------------------------------------------------------
INSERT INTO symptoms (code, name, description, category_id)
SELECT v.code, v.name, v.descr, cat.id
FROM (VALUES
  -- Neurological
  ('SYM-NUMB',   'Numbness',           'Tingling or numbness in hands and feet',   'Neurological'),
  ('SYM-TREM',   'Tremor',             'Involuntary shaking or trembling',         'Neurological'),
  ('SYM-PARA',   'Paralysis',          'Loss of muscle movement',                  'Neurological'),
  ('SYM-CONV',   'Convulsions',        'Seizures or violent spasms',               'Neurological'),
  ('SYM-DIZZ',   'Dizziness',          'Lightheadedness or vertigo',               'Neurological'),
  ('SYM-DELI',   'Delirium',           'Confusion and disordered thinking',        'Neurological'),
  ('SYM-HEAD',   'Headache',           'Persistent head pain',                     'Neurological'),
  ('SYM-WEAK',   'Muscle weakness',    'Reduced muscle strength',                  'Neurological'),
  -- Gastrointestinal
  ('SYM-NAUS',   'Nausea',             'A queasy, sick feeling',                   'Gastrointestinal'),
  ('SYM-VOMI',   'Vomiting',           'Forceful expulsion of stomach contents',   'Gastrointestinal'),
  ('SYM-ABDO',   'Abdominal pain',     'Cramping or pain in the belly',            'Gastrointestinal'),
  ('SYM-DIAR',   'Diarrhea',           'Loose or watery stools',                   'Gastrointestinal'),
  ('SYM-BLOOD',  'Bloody stool',       'Blood passed in the stool',                'Gastrointestinal'),
  ('SYM-APPE',   'Loss of appetite',   'Reduced desire to eat',                    'Gastrointestinal'),
  ('SYM-SALI',   'Excess salivation',  'Excessive drooling',                       'Gastrointestinal'),
  ('SYM-THIRST', 'Intense thirst',     'Unquenchable thirst',                      'Gastrointestinal'),
  -- Dermatological
  ('SYM-RASH',   'Rash',               'Skin eruption or spots',                   'Dermatological'),
  ('SYM-LEADLINE','Gum line',          'Blue-gray line along the gums',            'Dermatological'),
  ('SYM-JAUN',   'Jaundice',           'Yellowing of skin and eyes',               'Dermatological'),
  ('SYM-PALL',   'Pallor',             'Unusual paleness of the skin',             'Dermatological'),
  ('SYM-SWEAT',  'Sweating',           'Profuse sweating',                         'Dermatological'),
  ('SYM-ULCER',  'Ulcers',             'Open sores on skin or mouth',              'Dermatological'),
  ('SYM-PUST',   'Pustules',           'Fluid-filled blisters or pox',             'Dermatological'),
  ('SYM-BRUISE', 'Easy bruising',      'Bruising and bleeding gums',               'Dermatological'),
  -- Cardiovascular
  ('SYM-PALP',   'Palpitations',       'Racing or irregular heartbeat',            'Cardiovascular'),
  ('SYM-BRADY',  'Slow pulse',         'Abnormally slow heartbeat',                'Cardiovascular'),
  ('SYM-EDEMA',  'Edema',              'Swelling of legs and feet',                'Cardiovascular'),
  ('SYM-CHEST',  'Chest pain',         'Pain or tightness in the chest',           'Cardiovascular'),
  -- Respiratory
  ('SYM-DYSP',   'Shortness of breath','Difficulty breathing',                     'Respiratory'),
  ('SYM-COUGH',  'Cough',              'Persistent coughing',                      'Respiratory'),
  ('SYM-RESPPAR','Respiratory failure','Paralysis of the breathing muscles',       'Respiratory'),
  ('SYM-HEMOP',  'Coughing blood',     'Blood brought up when coughing',           'Respiratory'),
  -- Constitutional
  ('SYM-FEVER',  'Fever',              'Elevated body temperature',                'Constitutional'),
  ('SYM-FATI',   'Fatigue',            'Persistent tiredness',                     'Constitutional'),
  ('SYM-WEIGHT', 'Weight loss',        'Unexplained wasting',                      'Constitutional'),
  ('SYM-CHILL',  'Chills',             'Shivering and cold sweats',                'Constitutional'),
  -- Ophthalmic
  ('SYM-BLUR',   'Blurred vision',     'Loss of visual clarity',                   'Ophthalmic'),
  ('SYM-DILA',   'Dilated pupils',     'Abnormally widened pupils',                'Ophthalmic'),
  ('SYM-CONSTR', 'Pinpoint pupils',    'Abnormally constricted pupils',            'Ophthalmic'),
  -- Psychiatric
  ('SYM-HALL',   'Hallucinations',     'Seeing or hearing things unreal',          'Psychiatric'),
  ('SYM-REST',   'Restlessness',       'Agitation and anxiety',                    'Psychiatric'),
  ('SYM-INSOM',  'Insomnia',           'Inability to sleep',                       'Psychiatric')
) AS v(code, name, descr, category_name)
JOIN categories cat ON cat.name = v.category_name;

-- ----------------------------------------------------------------------------
-- Factors (contextual clues that shift suspicion)
-- ----------------------------------------------------------------------------
INSERT INTO factors (code, name, description) VALUES
  ('FAC-POWDER',      'Uses white face powder',   'Wears white lead cosmetic, as favored in the rear palace.'),
  ('FAC-RICE',        'Diet of polished rice',    'Eats refined white rice with the bran stripped away.'),
  ('FAC-PLEASURE',    'Frequents the pleasure district', 'Spends time in or works at the pleasure quarter.'),
  ('FAC-BANQUET',     'Attended a shared banquet', 'Ate and drank from a common court feast.'),
  ('FAC-REARPALACE',  'Resides in the rear palace', 'Lives within the inner palace of consorts.'),
  ('FAC-PREGNANT',    'Pregnant or recently gave birth', 'Currently with child or newly delivered.'),
  ('FAC-SUMMER',      'Warm humid season',        'Hot months of insects and swift spoilage.'),
  ('FAC-FORAGED',     'Ate foraged plants',       'Consumed wild-gathered plants, roots, or mushrooms.'),
  ('FAC-SEAFOOD',     'Ate fish or seafood',      'Recently ate fish, which may be improperly prepared.'),
  ('FAC-CROWDED',     'Crowded quarters',         'Lives or works among many people in close contact.'),
  ('FAC-ALCHEMY',     'Handles minerals',         'Works with cosmetics, pigments, or alchemical substances.'),
  ('FAC-WINE',        'Heavy drinking',           'Consumes large amounts of wine or strong liquor.'),
  ('FAC-HONEY',       'Consumed wild honey',      'Ate honey gathered from mountain flowers.'),
  ('FAC-SNAKEBITE',   'Recent bite',              'Was recently bitten by a snake or venomous creature.'),
  ('FAC-WELLWATER',   'Communal water source',    'Drinks from a shared well or standing water.'),
  ('FAC-APHRODISIAC', 'Given a love potion',      'Was administered an aphrodisiac or secret draught.');

-- ----------------------------------------------------------------------------
-- Causes (poisons + diseases); remedy_id resolved by remedy name
-- ----------------------------------------------------------------------------
INSERT INTO causes (code, name, type, prior, description, remedy_id)
SELECT v.code, v.name, v.ctype, v.prior, v.descr,
       (SELECT id FROM remedies r WHERE r.name = v.remedy_name)
FROM (VALUES
  -- POISONS -----------------------------------------------------------------
  ('POI-LEAD',  'Lead poisoning (white face powder)', 'poison', 0.050, 'Cumulative poisoning from white lead cosmetic.',            'Purge and Cease Powder'),
  ('POI-HG',    'Mercury poisoning',                  'poison', 0.020, 'Poisoning from quicksilver vapor or compounds.',            'Mercury Purge'),
  ('POI-ARS',   'Arsenic poisoning',                  'poison', 0.020, 'A tasteless mineral poison, the assassin''s favorite.',     'Activated Charcoal'),
  ('POI-ACON',  'Aconite (wolfsbane)',                'poison', 0.015, 'Toxin of the monkshood root; numbs then stops the heart.', 'Palace Physician''s Antidote'),
  ('POI-FUGU',  'Tetrodotoxin (pufferfish)',          'poison', 0.010, 'Nerve poison of the pufferfish; the mind stays clear.',     NULL),
  ('POI-CANT',  'Cantharidin (blister beetle)',       'poison', 0.010, 'Caustic aphrodisiac that burns the passages within.',      'Rehydration Broth'),
  ('POI-BELLA', 'Belladonna (nightshade)',            'poison', 0.010, 'Deadly nightshade; widens the eyes and unhinges the mind.','Palace Physician''s Antidote'),
  ('POI-DATU',  'Datura (thorn-apple)',               'poison', 0.008, 'Delirium-inducing seeds and flowers.',                     'Palace Physician''s Antidote'),
  ('POI-HONEY', 'Mad honey (grayanotoxin)',           'poison', 0.008, 'Honey from toxic blooms; slows the heart and head.',       'Palace Physician''s Antidote'),
  ('POI-SNAKE', 'Snake venom',                        'poison', 0.010, 'Venom from a bite, swelling and weakening the body.',       'Antivenom Serum'),
  ('POI-OPIUM', 'Opium overdose',                     'poison', 0.020, 'Excess of the poppy; breath and pulse sink into sleep.',    'Stimulant Draught'),
  ('POI-MUSH',  'Mushroom poisoning',                 'poison', 0.020, 'Toxins of misidentified mushrooms; the liver pays later.', 'Activated Charcoal'),
  ('POI-OLEA',  'Oleander poisoning',                 'poison', 0.008, 'Ornamental shrub whose every part strains the heart.',      'Palace Physician''s Antidote'),
  ('POI-STRY',  'Strychnine (nux vomica)',            'poison', 0.006, 'Convulsant seed poison; the body arches and locks.',        'Activated Charcoal'),
  ('POI-HEML',  'Hemlock poisoning',                  'poison', 0.006, 'Creeping paralysis rising from the legs to the breath.',    NULL),
  -- DISEASES ----------------------------------------------------------------
  ('DIS-BERI',  'Beriberi (thiamine deficiency)',     'disease', 0.120, 'The wasting sickness of those who eat only white rice.',   'Rice Bran Regimen'),
  ('DIS-POX',   'Smallpox',                           'disease', 0.030, 'Feared pox that scars and kills in waves.',                'Isolation and Rest'),
  ('DIS-MEAS',  'Measles',                            'disease', 0.040, 'Spreading rash and fever, dangerous to the young.',        'Isolation and Rest'),
  ('DIS-MALA',  'Malaria',                            'disease', 0.060, 'The shaking fever that returns in cycles.',                'Cinchona Bark'),
  ('DIS-TUBE',  'Consumption (tuberculosis)',         'disease', 0.040, 'A slow wasting cough that thins the patient away.',        'Isolation and Rest'),
  ('DIS-SYPH',  'The pox (syphilis)',                 'disease', 0.040, 'Illness of the pleasure quarter, from sores to ruin.',     'Mercury Purge'),
  ('DIS-TYPH',  'Typhoid fever',                      'disease', 0.050, 'Enteric fever from tainted food and water.',               'Rehydration Broth'),
  ('DIS-DYSE',  'Dysentery',                          'disease', 0.060, 'Bloody flux from filth in the water.',                     'Rehydration Broth'),
  ('DIS-SCUR',  'Scurvy',                             'disease', 0.030, 'Sickness of poor diet; the gums bleed and wounds gape.',   'Citrus and Greens'),
  ('DIS-FOOD',  'Food poisoning',                     'disease', 0.100, 'Acute upset from spoiled or tainted food.',                'Ginger Infusion'),
  ('DIS-CHOL',  'Cholera',                            'disease', 0.030, 'Violent flux that drains the body of water.',              'Rehydration Broth'),
  ('DIS-FLU',   'Influenza',                          'disease', 0.100, 'Seasonal fever, ache, and cough that sweeps a household.', 'Isolation and Rest'),
  ('DIS-ANEM',  'Anemia',                             'disease', 0.080, 'Thin, pale blood leaving the body tired and faint.',       'Liver and Iron Tonic'),
  ('DIS-PNEU',  'Pneumonia',                          'disease', 0.050, 'Deep chest fever with laboured, painful breathing.',       'Isolation and Rest')
) AS v(code, name, ctype, prior, descr, remedy_name);

-- ----------------------------------------------------------------------------
-- Symptom profiles: likelihood = P(symptom | cause)
-- ----------------------------------------------------------------------------
INSERT INTO symptom_profiles (cause_id, symptom_id, likelihood)
SELECT c.id, s.id, v.likelihood
FROM (VALUES
  -- Lead poisoning
  ('POI-LEAD','SYM-ABDO',0.70), ('POI-LEAD','SYM-LEADLINE',0.60), ('POI-LEAD','SYM-PALL',0.65),
  ('POI-LEAD','SYM-FATI',0.70), ('POI-LEAD','SYM-HEAD',0.55), ('POI-LEAD','SYM-NUMB',0.50),
  ('POI-LEAD','SYM-WEAK',0.55), ('POI-LEAD','SYM-CONV',0.30), ('POI-LEAD','SYM-WEIGHT',0.40),
  ('POI-LEAD','SYM-DELI',0.25),
  -- Mercury
  ('POI-HG','SYM-TREM',0.80), ('POI-HG','SYM-SALI',0.70), ('POI-HG','SYM-DELI',0.60),
  ('POI-HG','SYM-NUMB',0.50), ('POI-HG','SYM-INSOM',0.55), ('POI-HG','SYM-RASH',0.40),
  ('POI-HG','SYM-WEAK',0.45),
  -- Arsenic
  ('POI-ARS','SYM-VOMI',0.75), ('POI-ARS','SYM-DIAR',0.75), ('POI-ARS','SYM-ABDO',0.70),
  ('POI-ARS','SYM-THIRST',0.60), ('POI-ARS','SYM-WEAK',0.60), ('POI-ARS','SYM-RASH',0.30),
  ('POI-ARS','SYM-CONV',0.25),
  -- Aconite
  ('POI-ACON','SYM-NUMB',0.90), ('POI-ACON','SYM-PALP',0.75), ('POI-ACON','SYM-VOMI',0.55),
  ('POI-ACON','SYM-WEAK',0.60), ('POI-ACON','SYM-DIZZ',0.60), ('POI-ACON','SYM-SALI',0.50),
  ('POI-ACON','SYM-RESPPAR',0.40),
  -- Tetrodotoxin
  ('POI-FUGU','SYM-NUMB',0.90), ('POI-FUGU','SYM-PARA',0.80), ('POI-FUGU','SYM-WEAK',0.75),
  ('POI-FUGU','SYM-RESPPAR',0.70), ('POI-FUGU','SYM-DIZZ',0.50), ('POI-FUGU','SYM-NAUS',0.40),
  ('POI-FUGU','SYM-SALI',0.30),
  -- Cantharidin
  ('POI-CANT','SYM-ABDO',0.60), ('POI-CANT','SYM-BLOOD',0.40), ('POI-CANT','SYM-REST',0.55),
  ('POI-CANT','SYM-FEVER',0.40), ('POI-CANT','SYM-ULCER',0.40), ('POI-CANT','SYM-PALP',0.40),
  -- Belladonna
  ('POI-BELLA','SYM-DILA',0.90), ('POI-BELLA','SYM-BLUR',0.70), ('POI-BELLA','SYM-DELI',0.65),
  ('POI-BELLA','SYM-HALL',0.60), ('POI-BELLA','SYM-REST',0.60), ('POI-BELLA','SYM-PALP',0.60),
  ('POI-BELLA','SYM-FEVER',0.45),
  -- Datura
  ('POI-DATU','SYM-HALL',0.80), ('POI-DATU','SYM-DILA',0.85), ('POI-DATU','SYM-DELI',0.70),
  ('POI-DATU','SYM-REST',0.60), ('POI-DATU','SYM-BLUR',0.60), ('POI-DATU','SYM-FEVER',0.40),
  -- Mad honey
  ('POI-HONEY','SYM-DIZZ',0.80), ('POI-HONEY','SYM-BRADY',0.75), ('POI-HONEY','SYM-NAUS',0.60),
  ('POI-HONEY','SYM-VOMI',0.55), ('POI-HONEY','SYM-WEAK',0.60), ('POI-HONEY','SYM-BLUR',0.40),
  ('POI-HONEY','SYM-SWEAT',0.40),
  -- Snake venom
  ('POI-SNAKE','SYM-EDEMA',0.70), ('POI-SNAKE','SYM-PALL',0.50), ('POI-SNAKE','SYM-PALP',0.55),
  ('POI-SNAKE','SYM-DIZZ',0.55), ('POI-SNAKE','SYM-VOMI',0.45), ('POI-SNAKE','SYM-PARA',0.40),
  ('POI-SNAKE','SYM-BRUISE',0.50),
  -- Opium
  ('POI-OPIUM','SYM-CONSTR',0.90), ('POI-OPIUM','SYM-BRADY',0.60), ('POI-OPIUM','SYM-RESPPAR',0.55),
  ('POI-OPIUM','SYM-DIZZ',0.50), ('POI-OPIUM','SYM-NAUS',0.45), ('POI-OPIUM','SYM-DELI',0.30),
  -- Mushroom
  ('POI-MUSH','SYM-VOMI',0.80), ('POI-MUSH','SYM-DIAR',0.80), ('POI-MUSH','SYM-ABDO',0.75),
  ('POI-MUSH','SYM-JAUN',0.60), ('POI-MUSH','SYM-WEAK',0.55), ('POI-MUSH','SYM-DELI',0.30),
  -- Oleander
  ('POI-OLEA','SYM-PALP',0.70), ('POI-OLEA','SYM-BRADY',0.60), ('POI-OLEA','SYM-VOMI',0.60),
  ('POI-OLEA','SYM-NAUS',0.55), ('POI-OLEA','SYM-DIZZ',0.55), ('POI-OLEA','SYM-BLUR',0.40),
  -- Strychnine
  ('POI-STRY','SYM-CONV',0.90), ('POI-STRY','SYM-REST',0.60), ('POI-STRY','SYM-SWEAT',0.50),
  ('POI-STRY','SYM-PALP',0.50), ('POI-STRY','SYM-WEAK',0.40),
  -- Hemlock
  ('POI-HEML','SYM-PARA',0.75), ('POI-HEML','SYM-WEAK',0.75), ('POI-HEML','SYM-RESPPAR',0.60),
  ('POI-HEML','SYM-TREM',0.40), ('POI-HEML','SYM-SALI',0.40), ('POI-HEML','SYM-DIZZ',0.50),
  -- Beriberi
  ('DIS-BERI','SYM-EDEMA',0.75), ('DIS-BERI','SYM-NUMB',0.80), ('DIS-BERI','SYM-WEAK',0.75),
  ('DIS-BERI','SYM-FATI',0.80), ('DIS-BERI','SYM-PALP',0.60), ('DIS-BERI','SYM-WEIGHT',0.50),
  ('DIS-BERI','SYM-DYSP',0.50), ('DIS-BERI','SYM-APPE',0.55), ('DIS-BERI','SYM-DELI',0.30),
  -- Smallpox
  ('DIS-POX','SYM-PUST',0.90), ('DIS-POX','SYM-FEVER',0.90), ('DIS-POX','SYM-RASH',0.80),
  ('DIS-POX','SYM-HEAD',0.70), ('DIS-POX','SYM-FATI',0.70), ('DIS-POX','SYM-VOMI',0.40),
  -- Measles
  ('DIS-MEAS','SYM-RASH',0.90), ('DIS-MEAS','SYM-FEVER',0.85), ('DIS-MEAS','SYM-COUGH',0.75),
  ('DIS-MEAS','SYM-FATI',0.60), ('DIS-MEAS','SYM-DYSP',0.30),
  -- Malaria
  ('DIS-MALA','SYM-FEVER',0.90), ('DIS-MALA','SYM-CHILL',0.85), ('DIS-MALA','SYM-SWEAT',0.75),
  ('DIS-MALA','SYM-HEAD',0.60), ('DIS-MALA','SYM-FATI',0.70), ('DIS-MALA','SYM-PALL',0.50),
  ('DIS-MALA','SYM-VOMI',0.40),
  -- Consumption / TB
  ('DIS-TUBE','SYM-COUGH',0.85), ('DIS-TUBE','SYM-HEMOP',0.60), ('DIS-TUBE','SYM-WEIGHT',0.75),
  ('DIS-TUBE','SYM-FEVER',0.60), ('DIS-TUBE','SYM-FATI',0.75), ('DIS-TUBE','SYM-SWEAT',0.60),
  ('DIS-TUBE','SYM-DYSP',0.50),
  -- Syphilis
  ('DIS-SYPH','SYM-ULCER',0.70), ('DIS-SYPH','SYM-RASH',0.65), ('DIS-SYPH','SYM-FEVER',0.40),
  ('DIS-SYPH','SYM-FATI',0.50), ('DIS-SYPH','SYM-HEAD',0.40), ('DIS-SYPH','SYM-WEIGHT',0.30),
  -- Typhoid
  ('DIS-TYPH','SYM-FEVER',0.90), ('DIS-TYPH','SYM-ABDO',0.60), ('DIS-TYPH','SYM-DIAR',0.55),
  ('DIS-TYPH','SYM-HEAD',0.60), ('DIS-TYPH','SYM-FATI',0.70), ('DIS-TYPH','SYM-RASH',0.30),
  ('DIS-TYPH','SYM-DELI',0.35),
  -- Dysentery
  ('DIS-DYSE','SYM-DIAR',0.90), ('DIS-DYSE','SYM-BLOOD',0.80), ('DIS-DYSE','SYM-ABDO',0.75),
  ('DIS-DYSE','SYM-FEVER',0.60), ('DIS-DYSE','SYM-FATI',0.55), ('DIS-DYSE','SYM-THIRST',0.50),
  -- Scurvy
  ('DIS-SCUR','SYM-BRUISE',0.85), ('DIS-SCUR','SYM-FATI',0.70), ('DIS-SCUR','SYM-WEAK',0.65),
  ('DIS-SCUR','SYM-PALL',0.55), ('DIS-SCUR','SYM-ULCER',0.50), ('DIS-SCUR','SYM-WEIGHT',0.40),
  -- Food poisoning
  ('DIS-FOOD','SYM-VOMI',0.85), ('DIS-FOOD','SYM-DIAR',0.80), ('DIS-FOOD','SYM-ABDO',0.75),
  ('DIS-FOOD','SYM-NAUS',0.80), ('DIS-FOOD','SYM-FEVER',0.40), ('DIS-FOOD','SYM-WEAK',0.45),
  -- Cholera
  ('DIS-CHOL','SYM-DIAR',0.95), ('DIS-CHOL','SYM-VOMI',0.70), ('DIS-CHOL','SYM-THIRST',0.85),
  ('DIS-CHOL','SYM-WEAK',0.75), ('DIS-CHOL','SYM-DIZZ',0.55), ('DIS-CHOL','SYM-PALP',0.40),
  -- Influenza
  ('DIS-FLU','SYM-FEVER',0.85), ('DIS-FLU','SYM-COUGH',0.75), ('DIS-FLU','SYM-FATI',0.80),
  ('DIS-FLU','SYM-HEAD',0.70), ('DIS-FLU','SYM-CHILL',0.60), ('DIS-FLU','SYM-WEAK',0.55),
  -- Anemia
  ('DIS-ANEM','SYM-PALL',0.85), ('DIS-ANEM','SYM-FATI',0.85), ('DIS-ANEM','SYM-DIZZ',0.65),
  ('DIS-ANEM','SYM-PALP',0.55), ('DIS-ANEM','SYM-DYSP',0.50), ('DIS-ANEM','SYM-WEAK',0.60),
  -- Pneumonia
  ('DIS-PNEU','SYM-COUGH',0.85), ('DIS-PNEU','SYM-FEVER',0.85), ('DIS-PNEU','SYM-DYSP',0.75),
  ('DIS-PNEU','SYM-CHEST',0.60), ('DIS-PNEU','SYM-FATI',0.70), ('DIS-PNEU','SYM-CHILL',0.55)
) AS v(cause_code, symptom_code, likelihood)
JOIN causes   c ON c.code = v.cause_code
JOIN symptoms s ON s.code = v.symptom_code;

-- ----------------------------------------------------------------------------
-- Factor profiles: modifier = P(factor present | cause), must be <= 1
-- ----------------------------------------------------------------------------
INSERT INTO factor_profiles (cause_id, factor_id, modifier)
SELECT c.id, f.id, v.modifier
FROM (VALUES
  ('POI-LEAD','FAC-POWDER',0.90),   ('POI-LEAD','FAC-REARPALACE',0.65), ('POI-LEAD','FAC-PREGNANT',0.35), ('POI-LEAD','FAC-ALCHEMY',0.30),
  ('POI-HG','FAC-ALCHEMY',0.60),    ('POI-HG','FAC-POWDER',0.25),
  ('POI-ARS','FAC-ALCHEMY',0.45),   ('POI-ARS','FAC-BANQUET',0.40),
  ('POI-ACON','FAC-FORAGED',0.60),  ('POI-ACON','FAC-BANQUET',0.30),
  ('POI-FUGU','FAC-SEAFOOD',0.95),  ('POI-FUGU','FAC-BANQUET',0.30),
  ('POI-CANT','FAC-APHRODISIAC',0.85), ('POI-CANT','FAC-PLEASURE',0.40), ('POI-CANT','FAC-BANQUET',0.30), ('POI-CANT','FAC-WINE',0.30),
  ('POI-BELLA','FAC-FORAGED',0.70),
  ('POI-DATU','FAC-FORAGED',0.70),  ('POI-DATU','FAC-APHRODISIAC',0.30),
  ('POI-HONEY','FAC-HONEY',0.95),
  ('POI-SNAKE','FAC-SNAKEBITE',0.95),
  ('POI-OPIUM','FAC-PLEASURE',0.50),('POI-OPIUM','FAC-APHRODISIAC',0.20),
  ('POI-MUSH','FAC-FORAGED',0.95),
  ('POI-OLEA','FAC-FORAGED',0.55),
  ('POI-STRY','FAC-ALCHEMY',0.30),  ('POI-STRY','FAC-FORAGED',0.30),
  ('POI-HEML','FAC-FORAGED',0.70),
  ('DIS-BERI','FAC-RICE',0.95),     ('DIS-BERI','FAC-REARPALACE',0.50), ('DIS-BERI','FAC-PREGNANT',0.35), ('DIS-BERI','FAC-WINE',0.30),
  ('DIS-SYPH','FAC-PLEASURE',0.90),
  ('DIS-MALA','FAC-SUMMER',0.80),   ('DIS-MALA','FAC-WELLWATER',0.30),
  ('DIS-CHOL','FAC-WELLWATER',0.85),('DIS-CHOL','FAC-CROWDED',0.55),    ('DIS-CHOL','FAC-SUMMER',0.50),
  ('DIS-DYSE','FAC-WELLWATER',0.75),('DIS-DYSE','FAC-CROWDED',0.55),    ('DIS-DYSE','FAC-SUMMER',0.45),
  ('DIS-TYPH','FAC-WELLWATER',0.75),('DIS-TYPH','FAC-CROWDED',0.50),
  ('DIS-FOOD','FAC-BANQUET',0.55),  ('DIS-FOOD','FAC-SUMMER',0.45),     ('DIS-FOOD','FAC-SEAFOOD',0.35),
  ('DIS-POX','FAC-CROWDED',0.70),
  ('DIS-MEAS','FAC-CROWDED',0.70),
  ('DIS-FLU','FAC-CROWDED',0.65),
  ('DIS-TUBE','FAC-CROWDED',0.60),  ('DIS-TUBE','FAC-PLEASURE',0.30),
  ('DIS-SCUR','FAC-RICE',0.30),
  ('DIS-ANEM','FAC-PREGNANT',0.40), ('DIS-ANEM','FAC-REARPALACE',0.30), ('DIS-ANEM','FAC-RICE',0.30),
  ('DIS-PNEU','FAC-CROWDED',0.45)
) AS v(cause_code, factor_code, modifier)
JOIN causes  c ON c.code = v.cause_code
JOIN factors f ON f.code = v.factor_code;

-- ----------------------------------------------------------------------------
-- Cases (palace mysteries for Maomao to unravel)
-- ----------------------------------------------------------------------------
INSERT INTO cases (title, description, status) VALUES
  ('The Ailing Princes',        'Imperial infants sicken and fail to thrive within the rear palace.', 'closed'),
  ('Wasting in the Rear Palace','Consorts grow numb, swollen, and short of breath.',                  'closed'),
  ('The Poisoned Banquet',      'Guests collapse after a shared feast; the source is unclear.',       'open'),
  ('The Courtesan''s Malady',   'A celebrated courtesan is stricken with sores and fever.',           'open'),
  ('The Envoy''s Collapse',     'A foreign envoy falls gravely ill hours after a private meal.',      'open');

-- ----------------------------------------------------------------------------
-- Observations — symptom sightings (kind = 'symptom', ref_id -> symptoms.id)
-- ----------------------------------------------------------------------------
INSERT INTO observations (case_id, kind, ref_id, note)
SELECT ca.id, 'symptom', s.id, v.note
FROM (VALUES
  ('The Ailing Princes','SYM-PALL',    'Unusual paleness noted by the wet nurse.'),
  ('The Ailing Princes','SYM-ABDO',    'Persistent colic and crying.'),
  ('The Ailing Princes','SYM-LEADLINE','Faint dark line seen along the gums.'),
  ('The Ailing Princes','SYM-CONV',    'A brief fit witnessed by the physician.'),
  ('The Ailing Princes','SYM-FATI',    'Listless and slow to feed.'),
  ('Wasting in the Rear Palace','SYM-EDEMA','Swelling of the ankles and shins.'),
  ('Wasting in the Rear Palace','SYM-NUMB','Complaints of numb feet.'),
  ('Wasting in the Rear Palace','SYM-FATI','Constant weariness.'),
  ('Wasting in the Rear Palace','SYM-PALP','Fluttering heartbeat at rest.'),
  ('Wasting in the Rear Palace','SYM-APPE','Little appetite despite ample food.'),
  ('The Poisoned Banquet','SYM-NUMB',  'Tingling lips reported first.'),
  ('The Poisoned Banquet','SYM-PALP',  'Racing pulse among several guests.'),
  ('The Poisoned Banquet','SYM-VOMI',  'Sudden vomiting at the table.'),
  ('The Poisoned Banquet','SYM-DIZZ',  'Guests grew faint and unsteady.'),
  ('The Courtesan''s Malady','SYM-ULCER','A painless sore observed.'),
  ('The Courtesan''s Malady','SYM-RASH','Spreading rash across the trunk.'),
  ('The Courtesan''s Malady','SYM-FEVER','Low fever that will not break.'),
  ('The Courtesan''s Malady','SYM-FATI','Growing exhaustion.'),
  ('The Envoy''s Collapse','SYM-VOMI', 'Violent vomiting after the meal.'),
  ('The Envoy''s Collapse','SYM-DIAR', 'Severe watery flux.'),
  ('The Envoy''s Collapse','SYM-ABDO', 'Gripping abdominal pain.'),
  ('The Envoy''s Collapse','SYM-JAUN', 'Yellow tinge appearing in the eyes.')
) AS v(case_title, symptom_code, note)
JOIN cases    ca ON ca.title = v.case_title
JOIN symptoms s  ON s.code   = v.symptom_code;

-- ----------------------------------------------------------------------------
-- Observations — contextual factors (kind = 'factor', ref_id -> factors.id)
-- ----------------------------------------------------------------------------
INSERT INTO observations (case_id, kind, ref_id, note)
SELECT ca.id, 'factor', f.id, v.note
FROM (VALUES
  ('The Ailing Princes','FAC-POWDER',     'The nurse wears heavy white cosmetic.'),
  ('The Ailing Princes','FAC-REARPALACE', 'Confined to the inner palace since birth.'),
  ('The Ailing Princes','FAC-PREGNANT',   'The mother used the same powder while carrying.'),
  ('Wasting in the Rear Palace','FAC-RICE','Fed only polished white rice.'),
  ('Wasting in the Rear Palace','FAC-REARPALACE','All affected reside in the rear palace.'),
  ('The Poisoned Banquet','FAC-BANQUET',  'All ate from the same feast.'),
  ('The Poisoned Banquet','FAC-SEAFOOD',  'A fish dish was served.'),
  ('The Courtesan''s Malady','FAC-PLEASURE','Works in the pleasure quarter.'),
  ('The Courtesan''s Malady','FAC-APHRODISIAC','Received unknown draughts from a patron.'),
  ('The Envoy''s Collapse','FAC-FORAGED', 'The meal included gathered mountain fare.'),
  ('The Envoy''s Collapse','FAC-BANQUET', 'A private meal was shared beforehand.')
) AS v(case_title, factor_code, note)
JOIN cases   ca ON ca.title = v.case_title
JOIN factors f  ON f.code   = v.factor_code;

COMMIT;
