1. Overview of IDnow
   IDnow is an identity verification service that can handle multiple types of flows:

VideoIdent: A video-based identity verification session with a trained IDnow agent.
eID: Uses the German electronic ID card’s NFC chip for 100% digital identification.
eSign: Adds a qualified electronic signature (QES) after verifying the user’s identity (via video or eID).
InstantSign: Sometimes used if you already have user identity data on record, skipping a new video/eID check.
Depending on your user’s scenario, IDnow orchestrates:

Collecting user data (first name, last name, date of birth, etc.).
Potentially uploading or displaying contractual documents for eSigning.
A live or automated check of ID documents (passport, ID card).
Agent interaction or NFC-based eID reading for identity confirmation.
Generating or capturing a qualified signature on a PDF (for eSign).

2. Core Steps in an IDnow Identification
   Although the exact sequence may differ slightly by product, an IDnow identification typically looks like this:

Step 1: Create an Ident
Client/Partner (your application) calls IDnow’s POST /identifications/{transactionNumber}/start endpoint.
You provide user data (e.g. name, birthday, etc.) and specify a profile (e.g. “VIDEO_IDENT,” “EID,” or “ESIGN”).
IDnow responds with an internal ID (e.g. TST-XXXX) or transaction number confirmation.
The ident is now “created” or “started,” waiting for the user to proceed with the verification method.
Depending on the user’s channel (mobile app, web browser, or eID environment), IDnow can either:

Launch a VideoIdent session (the user connects in real-time with an IDnow agent).
Initiate an eID flow on a phone with an NFC reader.
Present documents for eSigning, or re-route them to a video or eID check if needed.
Step 2: Perform the Identification
2.1 VideoIdent
The user connects with an IDnow agent via live video (mobile or web).
The agent checks the user’s ID document for authenticity, verifies that the user’s face matches the document, and
possibly asks questions.
Once the agent is satisfied, the ident is set to finished or pending review (depending on your contract—some may go to
manual review).
2.2 eID
The user has a German electronic ID card and an NFC-capable phone.
IDnow’s SDK or app prompts the user to place their ID card on the phone’s NFC sensor.
The chip data is read (name, photo, etc.) and matched to the user.
If everything is valid, IDnow sets the ident to finished or verified.
2.3 eSign
If the user requires a Qualified Electronic Signature, the flow typically requires:
Identify the user (via video or eID).
Provide or upload the document(s) to be signed.
The user gives consent, and IDnow issues the QES via a recognized trust service provider.
The signed PDF is then stored or returned to your application with a digital signature certificate.
In some cases, for InstantSign, if you already have the user’s data in a compliant manner, IDnow might skip the
video/eID step and proceed directly to the sign step if IDnow’s rules are satisfied.

Step 3: Post-Verification
IDnow might run an internal or manual review.
The ident receives a final status: “finished,” “archived,” “canceled,” or “failed.”
You can retrieve the ident status or details by calling GET /identifications/{transactionNumber} or checking the final
status in your callback/webhook (if configured).
Step 4: Archiving or Deleting
Once done, you may archive or delete the ident from IDnow’s side (if your contract or data retention policy allows).
For archiving, you typically call POST /identifications/{transactionNumber}/archive.
IDnow often automatically archives after a certain period if the user does not complete the process.

3. IDnow eID vs. VideoIdent vs. eSign Flows
   VideoIdent: Real-time agent, valid in many compliance contexts (like AML in finance). The user sees a “live chat” or
   “video call” screen, and the agent checks the ID.
   eID: 100% digital, no agent involved, but restricted to the German ID card with NFC. Users must have ID card and a
   smartphone with NFC (and IDnow’s eID-compatible app).
   eSign: Involves a QES signature after verification. The QES issuance is only permitted once the user is verified;
   IDnow frequently does a short “video step” or uses eID.
   InstantSign: If you (the partner) already have user identity data, IDnow might skip a new ident and proceed with an
   immediate signature if IDnow’s policies allow.

6. Technical Summary of Typical Calls
   Create Ident
   POST /api/v1/{customer}/identifications/{transactionNumber}/start
   Parameters: user personal data (name, address, birthday, etc.), plus profile = VIDEO_IDENT | EID | ESIGN, etc.
   Response: JSON with internal ID or success.
   User completes the verification (video or eID reading). If eSign, user signs PDF(s).
   Check status:
   GET /api/v1/{customer}/identifications/{transactionNumber}
   -Response: JSON stating the status: “pending,” “started,” “finished,” “canceled,” etc.
   Archive ident (optional):
   POST /api/v1/{customer}/identifications/{transactionNumber}/archive
   No body on success (2xx).
   Additional endpoints might include:
   POST /api/v1/{customer}/documentdefinitions for storing repeated documents.
   POST /api/v1/{customer}/document/{uid}/check or POST /api/v1/{customer}/document/check for custom doc checks.
   GET /api/v1/{customer}/identifications/{transactionNumber}/report if your plan includes a PDF “report” (sometimes
   restricted by contract).
5. Key Differences in Each Flow
   5.1 VideoIdent
   Human agent or “testbot” for QA or test scenarios.
   The user must have a working camera, mic, stable connection, plus their physical ID.
   The agent guides them, capturing ID images.
   5.2 eID
   No agent required.
   Requires German ID card with RFID chip.
   The user must have an NFC-capable phone and the IDnow eID app or SDK.
   The system reads the card’s chip data, compares it to the user’s self-entered data, and returns a status.
   5.3 eSign
   Legal: Issues a Qualified Electronic Signature based on the identification.
   Typically uses a trust service provider in the background (like Namirial or D-TRUST).
   You can optionally upload documents to be signed. IDnow merges the QES once the user is verified.
6. Typical Developer Steps
   Obtain credentials (X-API-LOGIN-TOKEN, BaseURL, customer).
   Create an ident (POST .../start) specifying which product flow you want (VideoIdent, eID, eSign).
   Poll or listen to see if the ident was started, finished, or canceled.
   If eSign, upload or define the documents to sign. IDnow merges the signature after the user is confirmed.
   Archive or delete if desired, or fetch final results for your records.
7. Summary
   IDnow can handle different identity verification modes.
   Your system calls IDnow’s REST endpoints to create an ident, specifying user data and flow type.
   The user completes the verification (VideoIdent, eID, or eSign).
   You check status, optionally archive or retrieve data.
   Hence, the overall workflow is:

Set up IDnow credentials & environment (sandbox or prod).
POST “create ident” with user data, specifying VideoIdent, eID, or eSign.
Wait for user to finish or see if user canceled or timed out.
Retrieve ident details, or archive ident when done.
That’s a simplified explanation of the IDnow ident life cycle—covering eID, VideoIdent, and eSign in a single integrated
platform.