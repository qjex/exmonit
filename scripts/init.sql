CREATE TABLE public.rate (
                                           pair character varying(10) NOT NULL,
                                           exchange character varying(15) NOT NULL,
                                           rate character varying(30) NOT NULL,
                                           updated timestamp without time zone NOT NULL
);

ALTER TABLE public.rate OWNER TO postgres;

ALTER TABLE ONLY public.rate
    ADD CONSTRAINT rate_pk PRIMARY KEY (pair, exchange);