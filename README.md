# manager
💵 Automated extraction and insertion of credit card statement items into my Notion database

- Utilizes the bank's API to efficiently extract credit card statement data
- Filters transactions based on customizable time ranges (e.g., today, yesterday to today, last two days to today, etc.)
- Enhances transaction clarity by mapping transaction names from the API response using a configurable mapper config file
- Streamlines categorization by mapping transaction categories based on the mapped name, with fallback to a default category mapping if necessary
- Provides a transparent confirmation of selected transactions before insertion into the database
- Automatically adds new transactions (those not yet in the database) to your Notion database, ensuring comprehensive tracking and management.
