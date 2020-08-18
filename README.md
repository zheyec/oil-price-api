# oil-prices

This is an oil price look-up API. It was originally made for a chatbot product that was a real-world business project.

## Data Source

http://www.qiyoujiage.com/

## Product Requirements

### Format

Single-round task-based chatbot

### Triggers

{Province}{Gas Type} Price

Note:

- {Province}: must be provided. If not, alert "Please input the right province."
- {Gas Type}: optional. Includes #92, #95, #98, and #0, and their synonyms.

Example:

> #92 Gas Price in Shanghai
>
> Jiangsu Oil Price
>
> Beijing Diesel Price

### Chatbot Response

- If gas type is provided: Price of {Gas Type} today in {Province}: {Price}

- If gas type is not provided:

  - Oil price today in {Province}:

    "#92 Gas": {Price},
		"#95 Gas": {Price},
		"#98 Gas": {Price},
		"#0 Diesel": {Price}.

Note:

- Gas Type: use "#92 Gas", "#95 Gas", "#98 Gas", "#0 Diesel" as the standard wording.
- If price is unavailable, return NA.

## Local Test

```bash
./run.sh
```
