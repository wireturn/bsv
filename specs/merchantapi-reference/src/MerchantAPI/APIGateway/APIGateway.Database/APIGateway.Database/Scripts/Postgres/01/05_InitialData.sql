-- Copyright (c) 2020 Bitcoin Association.
-- Distributed under the Open BSV software license, see the accompanying file LICENSE

-- insert initial data into FeeQuote

do $$
declare
  cnt integer;
  feeQuoteId integer;
  standardFeeId integer;
  dataFeeId integer;
begin
  SELECT count(*) INTO cnt FROM public.feequote;
  if cnt = 0 then
    INSERT INTO public.feequote(createdat, validfrom, identity, identityprovider) 
    VALUES (now() at time zone 'utc', now() at time zone 'utc', null, null) 
    returning id INTO feeQuoteId;

    INSERT INTO public.Fee(feeQuote, feeType)
    VALUES(feeQuoteId, 'standard')
    returning id INTO standardFeeId;
	
	INSERT INTO public.FeeAmount (fee, satoshis, bytes) 
	VALUES(standardFeeId, 100, 200);
	INSERT INTO public.FeeAmount (fee, satoshis, bytes) 
	VALUES(standardFeeId, 100, 200);
	
	INSERT INTO public.Fee(feeQuote, feeType)
    VALUES(feeQuoteId, 'data')
    returning id INTO dataFeeId;
	
	INSERT INTO public.FeeAmount (fee, satoshis, bytes) 
	VALUES(dataFeeId, 100, 200);
	INSERT INTO public.FeeAmount (fee, satoshis, bytes) 
	VALUES(dataFeeId, 100, 200);
  end if;
end $$;
