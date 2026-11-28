--
-- PostgreSQL database dump
--

-- Dumped from database version 12.22 (Ubuntu 12.22-0ubuntu0.20.04.4)
-- Dumped by pg_dump version 12.22 (Ubuntu 12.22-0ubuntu0.20.04.4)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

DROP DATABASE universe;
--
-- Name: universe; Type: DATABASE; Schema: -; Owner: freecodecamp
--

CREATE DATABASE universe WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'C.UTF-8' LC_CTYPE = 'C.UTF-8';


ALTER DATABASE universe OWNER TO freecodecamp;

\connect universe

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: galaxy; Type: TABLE; Schema: public; Owner: freecodecamp
--

CREATE TABLE public.galaxy (
    galaxy_id integer NOT NULL,
    name character varying(225),
    decimal_age numeric(4,1),
    description text,
    descriptionxxx text,
    nn integer DEFAULT 1 NOT NULL,
    nnn integer DEFAULT 1 NOT NULL,
    un integer NOT NULL
);


ALTER TABLE public.galaxy OWNER TO freecodecamp;

--
-- Name: galaxy_id_seq; Type: SEQUENCE; Schema: public; Owner: freecodecamp
--

CREATE SEQUENCE public.galaxy_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.galaxy_id_seq OWNER TO freecodecamp;

--
-- Name: galaxy_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: freecodecamp
--

ALTER SEQUENCE public.galaxy_id_seq OWNED BY public.galaxy.galaxy_id;


--
-- Name: galaxy_un_seq; Type: SEQUENCE; Schema: public; Owner: freecodecamp
--

CREATE SEQUENCE public.galaxy_un_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.galaxy_un_seq OWNER TO freecodecamp;

--
-- Name: galaxy_un_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: freecodecamp
--

ALTER SEQUENCE public.galaxy_un_seq OWNED BY public.galaxy.un;


--
-- Name: moon; Type: TABLE; Schema: public; Owner: freecodecamp
--

CREATE TABLE public.moon (
    moon_id integer NOT NULL,
    name character varying(225),
    planet_id integer,
    description text,
    descriptionxxx text,
    nnn integer DEFAULT 1 NOT NULL,
    nn integer DEFAULT 1 NOT NULL,
    un integer NOT NULL
);


ALTER TABLE public.moon OWNER TO freecodecamp;

--
-- Name: moon_id_seq; Type: SEQUENCE; Schema: public; Owner: freecodecamp
--

CREATE SEQUENCE public.moon_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.moon_id_seq OWNER TO freecodecamp;

--
-- Name: moon_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: freecodecamp
--

ALTER SEQUENCE public.moon_id_seq OWNED BY public.moon.moon_id;


--
-- Name: moon_un_seq; Type: SEQUENCE; Schema: public; Owner: freecodecamp
--

CREATE SEQUENCE public.moon_un_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.moon_un_seq OWNER TO freecodecamp;

--
-- Name: moon_un_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: freecodecamp
--

ALTER SEQUENCE public.moon_un_seq OWNED BY public.moon.un;


--
-- Name: planet; Type: TABLE; Schema: public; Owner: freecodecamp
--

CREATE TABLE public.planet (
    planet_id integer NOT NULL,
    name character varying(225),
    age_in_millions_of_years integer,
    distance_from_earth integer,
    description text,
    is_spherical boolean,
    is_close_to_earth boolean,
    star_id integer,
    descriptionxxx text,
    nn integer DEFAULT 1 NOT NULL,
    nnn integer DEFAULT 1 NOT NULL,
    un integer NOT NULL
);


ALTER TABLE public.planet OWNER TO freecodecamp;

--
-- Name: planet_id_seq; Type: SEQUENCE; Schema: public; Owner: freecodecamp
--

CREATE SEQUENCE public.planet_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.planet_id_seq OWNER TO freecodecamp;

--
-- Name: planet_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: freecodecamp
--

ALTER SEQUENCE public.planet_id_seq OWNED BY public.planet.planet_id;


--
-- Name: planet_un_seq; Type: SEQUENCE; Schema: public; Owner: freecodecamp
--

CREATE SEQUENCE public.planet_un_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.planet_un_seq OWNER TO freecodecamp;

--
-- Name: planet_un_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: freecodecamp
--

ALTER SEQUENCE public.planet_un_seq OWNED BY public.planet.un;


--
-- Name: specification; Type: TABLE; Schema: public; Owner: freecodecamp
--

CREATE TABLE public.specification (
    specification_id integer NOT NULL,
    name character varying(225),
    color character varying(50),
    un integer NOT NULL
);


ALTER TABLE public.specification OWNER TO freecodecamp;

--
-- Name: specification_id_seq; Type: SEQUENCE; Schema: public; Owner: freecodecamp
--

CREATE SEQUENCE public.specification_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.specification_id_seq OWNER TO freecodecamp;

--
-- Name: specification_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: freecodecamp
--

ALTER SEQUENCE public.specification_id_seq OWNED BY public.specification.specification_id;


--
-- Name: specification_un_seq; Type: SEQUENCE; Schema: public; Owner: freecodecamp
--

CREATE SEQUENCE public.specification_un_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.specification_un_seq OWNER TO freecodecamp;

--
-- Name: specification_un_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: freecodecamp
--

ALTER SEQUENCE public.specification_un_seq OWNED BY public.specification.un;


--
-- Name: star; Type: TABLE; Schema: public; Owner: freecodecamp
--

CREATE TABLE public.star (
    star_id integer NOT NULL,
    name character varying(225),
    galaxy_id integer,
    description text,
    descriptionxxx text,
    nnn integer DEFAULT 1 NOT NULL,
    nn integer DEFAULT 1 NOT NULL,
    un integer NOT NULL
);


ALTER TABLE public.star OWNER TO freecodecamp;

--
-- Name: star_id_seq; Type: SEQUENCE; Schema: public; Owner: freecodecamp
--

CREATE SEQUENCE public.star_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.star_id_seq OWNER TO freecodecamp;

--
-- Name: star_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: freecodecamp
--

ALTER SEQUENCE public.star_id_seq OWNED BY public.star.star_id;


--
-- Name: star_un_seq; Type: SEQUENCE; Schema: public; Owner: freecodecamp
--

CREATE SEQUENCE public.star_un_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.star_un_seq OWNER TO freecodecamp;

--
-- Name: star_un_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: freecodecamp
--

ALTER SEQUENCE public.star_un_seq OWNED BY public.star.un;


--
-- Name: galaxy galaxy_id; Type: DEFAULT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.galaxy ALTER COLUMN galaxy_id SET DEFAULT nextval('public.galaxy_id_seq'::regclass);


--
-- Name: galaxy un; Type: DEFAULT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.galaxy ALTER COLUMN un SET DEFAULT nextval('public.galaxy_un_seq'::regclass);


--
-- Name: moon moon_id; Type: DEFAULT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.moon ALTER COLUMN moon_id SET DEFAULT nextval('public.moon_id_seq'::regclass);


--
-- Name: moon un; Type: DEFAULT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.moon ALTER COLUMN un SET DEFAULT nextval('public.moon_un_seq'::regclass);


--
-- Name: planet planet_id; Type: DEFAULT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.planet ALTER COLUMN planet_id SET DEFAULT nextval('public.planet_id_seq'::regclass);


--
-- Name: planet un; Type: DEFAULT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.planet ALTER COLUMN un SET DEFAULT nextval('public.planet_un_seq'::regclass);


--
-- Name: specification specification_id; Type: DEFAULT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.specification ALTER COLUMN specification_id SET DEFAULT nextval('public.specification_id_seq'::regclass);


--
-- Name: specification un; Type: DEFAULT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.specification ALTER COLUMN un SET DEFAULT nextval('public.specification_un_seq'::regclass);


--
-- Name: star star_id; Type: DEFAULT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.star ALTER COLUMN star_id SET DEFAULT nextval('public.star_id_seq'::regclass);


--
-- Name: star un; Type: DEFAULT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.star ALTER COLUMN un SET DEFAULT nextval('public.star_un_seq'::regclass);


--
-- Data for Name: galaxy; Type: TABLE DATA; Schema: public; Owner: freecodecamp
--

INSERT INTO public.galaxy VALUES (1, 'a', 3.0, NULL, NULL, 1, 1, 1);
INSERT INTO public.galaxy VALUES (2, 'aa', 3.0, NULL, NULL, 1, 1, 2);
INSERT INTO public.galaxy VALUES (3, 'ase', 3.0, NULL, NULL, 1, 1, 3);
INSERT INTO public.galaxy VALUES (4, 'asesaeae', 3.0, NULL, NULL, 1, 1, 4);
INSERT INTO public.galaxy VALUES (5, 'asesaeaseae', 3.0, NULL, NULL, 1, 1, 5);
INSERT INTO public.galaxy VALUES (6, 'asesaeaseaesssss', 3.0, NULL, NULL, 1, 1, 6);
INSERT INTO public.galaxy VALUES (7, 'asesaeaseaesssssxxxxxxxxx', 3.0, NULL, NULL, 1, 1, 7);


--
-- Data for Name: moon; Type: TABLE DATA; Schema: public; Owner: freecodecamp
--

INSERT INTO public.moon VALUES (1, 'x', 1, NULL, NULL, 1, 1, 1);
INSERT INTO public.moon VALUES (2, 'x', 2, NULL, NULL, 1, 1, 2);
INSERT INTO public.moon VALUES (3, 'x', 3, NULL, NULL, 1, 1, 3);
INSERT INTO public.moon VALUES (4, 'x', 4, NULL, NULL, 1, 1, 4);
INSERT INTO public.moon VALUES (5, 'x', 5, NULL, NULL, 1, 1, 5);
INSERT INTO public.moon VALUES (6, 'x', 6, NULL, NULL, 1, 1, 6);
INSERT INTO public.moon VALUES (7, 'x', 7, NULL, NULL, 1, 1, 7);
INSERT INTO public.moon VALUES (8, 'x', 8, NULL, NULL, 1, 1, 8);
INSERT INTO public.moon VALUES (9, 'x', 9, NULL, NULL, 1, 1, 9);
INSERT INTO public.moon VALUES (10, 'x', 10, NULL, NULL, 1, 1, 10);
INSERT INTO public.moon VALUES (11, 'x', 11, NULL, NULL, 1, 1, 11);
INSERT INTO public.moon VALUES (12, 'x', 12, NULL, NULL, 1, 1, 12);
INSERT INTO public.moon VALUES (13, 'x', 13, NULL, NULL, 1, 1, 13);
INSERT INTO public.moon VALUES (14, 'x', 14, NULL, NULL, 1, 1, 14);
INSERT INTO public.moon VALUES (16, 'x', 14, NULL, NULL, 1, 1, 15);
INSERT INTO public.moon VALUES (17, 'x', 14, NULL, NULL, 1, 1, 16);
INSERT INTO public.moon VALUES (18, 'x', 14, NULL, NULL, 1, 1, 17);
INSERT INTO public.moon VALUES (19, 'x', 14, NULL, NULL, 1, 1, 18);
INSERT INTO public.moon VALUES (20, 'x', 14, NULL, NULL, 1, 1, 19);
INSERT INTO public.moon VALUES (21, 'x', 14, NULL, NULL, 1, 1, 20);
INSERT INTO public.moon VALUES (22, 'x', 14, NULL, NULL, 1, 1, 21);
INSERT INTO public.moon VALUES (23, 'x', 14, NULL, NULL, 1, 1, 22);
INSERT INTO public.moon VALUES (24, 'x', 14, NULL, NULL, 1, 1, 23);
INSERT INTO public.moon VALUES (25, 'x', 14, NULL, NULL, 1, 1, 24);
INSERT INTO public.moon VALUES (26, 'x', 14, NULL, NULL, 1, 1, 25);


--
-- Data for Name: planet; Type: TABLE DATA; Schema: public; Owner: freecodecamp
--

INSERT INTO public.planet VALUES (1, 'x', NULL, NULL, NULL, NULL, NULL, 1, NULL, 1, 1, 1);
INSERT INTO public.planet VALUES (2, 'x', NULL, NULL, NULL, NULL, NULL, 2, NULL, 1, 1, 2);
INSERT INTO public.planet VALUES (3, 'x', NULL, NULL, NULL, NULL, NULL, 3, NULL, 1, 1, 3);
INSERT INTO public.planet VALUES (4, 'x', NULL, NULL, NULL, NULL, NULL, 4, NULL, 1, 1, 4);
INSERT INTO public.planet VALUES (5, 'x', NULL, NULL, NULL, NULL, NULL, 5, NULL, 1, 1, 5);
INSERT INTO public.planet VALUES (6, 'x', NULL, NULL, NULL, NULL, NULL, 6, NULL, 1, 1, 6);
INSERT INTO public.planet VALUES (7, 'xx', NULL, NULL, NULL, NULL, NULL, 6, NULL, 1, 1, 7);
INSERT INTO public.planet VALUES (8, 'xxx', NULL, NULL, NULL, NULL, NULL, 6, NULL, 1, 1, 8);
INSERT INTO public.planet VALUES (9, 'xxxx', NULL, NULL, NULL, NULL, NULL, 6, NULL, 1, 1, 9);
INSERT INTO public.planet VALUES (10, 'xxxxxx', NULL, NULL, NULL, NULL, NULL, 6, NULL, 1, 1, 10);
INSERT INTO public.planet VALUES (11, 'xxxxxx123', NULL, NULL, NULL, NULL, NULL, 6, NULL, 1, 1, 11);
INSERT INTO public.planet VALUES (12, 'xxxxxx11123', NULL, NULL, NULL, NULL, NULL, 6, NULL, 1, 1, 12);
INSERT INTO public.planet VALUES (13, 'xxxxxx1112xxx3', NULL, NULL, NULL, NULL, NULL, 6, NULL, 1, 1, 13);
INSERT INTO public.planet VALUES (14, 'xxxxxx1112xxxxx3', NULL, NULL, NULL, NULL, NULL, 6, NULL, 1, 1, 14);


--
-- Data for Name: specification; Type: TABLE DATA; Schema: public; Owner: freecodecamp
--

INSERT INTO public.specification VALUES (1, 'x', NULL, 1);
INSERT INTO public.specification VALUES (2, 'x', NULL, 2);
INSERT INTO public.specification VALUES (3, 'x', NULL, 3);
INSERT INTO public.specification VALUES (4, 'x', NULL, 4);
INSERT INTO public.specification VALUES (5, 'x', NULL, 5);


--
-- Data for Name: star; Type: TABLE DATA; Schema: public; Owner: freecodecamp
--

INSERT INTO public.star VALUES (1, 'x', 1, NULL, NULL, 1, 1, 1);
INSERT INTO public.star VALUES (2, 'xxx', 2, NULL, NULL, 1, 1, 2);
INSERT INTO public.star VALUES (3, 'xxxaa', 3, NULL, NULL, 1, 1, 3);
INSERT INTO public.star VALUES (4, 'xxxaaaesae', 4, NULL, NULL, 1, 1, 4);
INSERT INTO public.star VALUES (5, 'xxxaaaesaeeeeeee', 5, NULL, NULL, 1, 1, 5);
INSERT INTO public.star VALUES (6, 'xxxaaaesaeeeeeeessssssssss', 6, NULL, NULL, 1, 1, 6);


--
-- Name: galaxy_id_seq; Type: SEQUENCE SET; Schema: public; Owner: freecodecamp
--

SELECT pg_catalog.setval('public.galaxy_id_seq', 7, true);


--
-- Name: galaxy_un_seq; Type: SEQUENCE SET; Schema: public; Owner: freecodecamp
--

SELECT pg_catalog.setval('public.galaxy_un_seq', 7, true);


--
-- Name: moon_id_seq; Type: SEQUENCE SET; Schema: public; Owner: freecodecamp
--

SELECT pg_catalog.setval('public.moon_id_seq', 26, true);


--
-- Name: moon_un_seq; Type: SEQUENCE SET; Schema: public; Owner: freecodecamp
--

SELECT pg_catalog.setval('public.moon_un_seq', 25, true);


--
-- Name: planet_id_seq; Type: SEQUENCE SET; Schema: public; Owner: freecodecamp
--

SELECT pg_catalog.setval('public.planet_id_seq', 14, true);


--
-- Name: planet_un_seq; Type: SEQUENCE SET; Schema: public; Owner: freecodecamp
--

SELECT pg_catalog.setval('public.planet_un_seq', 14, true);


--
-- Name: specification_id_seq; Type: SEQUENCE SET; Schema: public; Owner: freecodecamp
--

SELECT pg_catalog.setval('public.specification_id_seq', 5, true);


--
-- Name: specification_un_seq; Type: SEQUENCE SET; Schema: public; Owner: freecodecamp
--

SELECT pg_catalog.setval('public.specification_un_seq', 5, true);


--
-- Name: star_id_seq; Type: SEQUENCE SET; Schema: public; Owner: freecodecamp
--

SELECT pg_catalog.setval('public.star_id_seq', 6, true);


--
-- Name: star_un_seq; Type: SEQUENCE SET; Schema: public; Owner: freecodecamp
--

SELECT pg_catalog.setval('public.star_un_seq', 6, true);


--
-- Name: galaxy galaxy_pkey; Type: CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.galaxy
    ADD CONSTRAINT galaxy_pkey PRIMARY KEY (galaxy_id);


--
-- Name: galaxy galaxy_un_key; Type: CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.galaxy
    ADD CONSTRAINT galaxy_un_key UNIQUE (un);


--
-- Name: moon moon_pkey; Type: CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.moon
    ADD CONSTRAINT moon_pkey PRIMARY KEY (moon_id);


--
-- Name: moon moon_un_key; Type: CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.moon
    ADD CONSTRAINT moon_un_key UNIQUE (un);


--
-- Name: planet planet_pkey; Type: CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.planet
    ADD CONSTRAINT planet_pkey PRIMARY KEY (planet_id);


--
-- Name: planet planet_un_key; Type: CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.planet
    ADD CONSTRAINT planet_un_key UNIQUE (un);


--
-- Name: specification specification_pkey; Type: CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.specification
    ADD CONSTRAINT specification_pkey PRIMARY KEY (specification_id);


--
-- Name: specification specification_un_key; Type: CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.specification
    ADD CONSTRAINT specification_un_key UNIQUE (un);


--
-- Name: star star_galaxy_id_key; Type: CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.star
    ADD CONSTRAINT star_galaxy_id_key UNIQUE (galaxy_id);


--
-- Name: star star_pkey; Type: CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.star
    ADD CONSTRAINT star_pkey PRIMARY KEY (star_id);


--
-- Name: star star_un_key; Type: CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.star
    ADD CONSTRAINT star_un_key UNIQUE (un);


--
-- Name: moon moon_planet_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.moon
    ADD CONSTRAINT moon_planet_id_fkey FOREIGN KEY (planet_id) REFERENCES public.planet(planet_id);


--
-- Name: planet planet_star_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.planet
    ADD CONSTRAINT planet_star_id_fkey FOREIGN KEY (star_id) REFERENCES public.star(star_id);


--
-- Name: star star_galaxy_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: freecodecamp
--

ALTER TABLE ONLY public.star
    ADD CONSTRAINT star_galaxy_id_fkey FOREIGN KEY (galaxy_id) REFERENCES public.galaxy(galaxy_id);


--
-- PostgreSQL database dump complete
--

