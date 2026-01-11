// Final Edition : This code add a course with crn code, then try to register it. then make submit
// The next step is to replace the register with register and drop with conditional add and drop

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"os/exec"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
)

func main() {
	available, err := isAvailable()
	if err != nil {
		fmt.Println("Error checking availability:", err)
		return
	}

	if available {
		fmt.Println("Course is available! Attempting to register...")
		// Add logic to Trigger the "Submit" request here if needed
		sendDesktopAlert()
		// 3. Email (Requires App Password)
		sendEmail()
	} else {
		fmt.Println("Course is not available / Closed Section. Retrying...")
		time.Sleep(60 * time.Second)
		main()
	}
}

func isAvailable() (bool, error) {
	// --- REQUEST 1: CHECK AVAILABILITY ---
	// Note: I removed the static 'Content-Length' header from the string so Go calculates it fresh.
	RawRequest1 := `POST /StudentRegistrationSsb/ssb/classRegistration/addCRNRegistrationItems HTTP/1.1
Host: banner9-registration.kfupm.edu.sa
Cookie: JSESSIONID=760032F45B31FD0370D2B9FAE3923EEE; _ga=GA1.3.1500780047.1767691760; _gid=GA1.3.2012430359.1767691760; _ga_EM00LC24FH=GS2.3.s1767691761$o1$g1$t1767691772$j49$l0$h0; KFUPM_Cookie=!F+vriIDI35PUtueAG8WAw/NalnvuQcHZoIEJ5QmvFhbsAFC9keeTENAtq/2tKv/e0XwdyTYNYGr9yLA=
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:146.0) Gecko/20100101 Firefox/146.0
Accept: application/json, text/javascript, */*; q=0.01
Accept-Language: en-US,en;q=0.5
Accept-Encoding: gzip, deflate, br
X-Synchronizer-Token: 3a81e06f-805f-47b9-a490-21156f23400c
Content-Type: application/x-www-form-urlencoded; charset=UTF-8
X-Requested-With: XMLHttpRequest
Content-Length: 25
Origin: https://banner9-registration.kfupm.edu.sa
Referer: https://banner9-registration.kfupm.edu.sa/StudentRegistrationSsb/ssb/classRegistration/classRegistration
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: same-origin
Priority: u=0
Te: trailers
Connection: keep-alive

crnList=22729&term=202520`

	req1, err := parseRawRequest(RawRequest1)
	if err != nil {
		return false, fmt.Errorf("failed to parse Req1: %v", err)
	}

	// EXECUTE REQUEST 1
	client := &http.Client{}
	resp, err := client.Do(req1)
	if err != nil {
		return false, fmt.Errorf("req1 failed: %v", err)
	}
	defer resp.Body.Close()

	// Debug: Check if Request 1 actually worked
	// body1, _ := io.ReadAll(resp.Body)
	// fmt.Println("Req1 Status:", resp.Status)

	time.Sleep(2 * time.Second)

	// --- REQUEST 2: SUBMIT / VERIFY ---
	// Note: I removed the Content-Length header here too.
	RawRequest2 := `POST /StudentRegistrationSsb/ssb/classRegistration/submitRegistration/batch HTTP/1.1
Host: banner9-registration.kfupm.edu.sa
Cookie: JSESSIONID=760032F45B31FD0370D2B9FAE3923EEE; _ga=GA1.3.1500780047.1767691760; _gid=GA1.3.2012430359.1767691760; _ga_EM00LC24FH=GS2.3.s1767691761$o1$g1$t1767691772$j49$l0$h0; KFUPM_Cookie=!F+vriIDI35PUtueAG8WAw/NalnvuQcHZoIEJ5QmvFhbsAFC9keeTENAtq/2tKv/e0XwdyTYNYGr9yLA=
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:146.0) Gecko/20100101 Firefox/146.0
Accept: application/json, text/javascript, */*; q=0.01
Accept-Language: en-US,en;q=0.5
Accept-Encoding: gzip, deflate, br
X-Synchronizer-Token: 3a81e06f-805f-47b9-a490-21156f23400c
Content-Type: application/json
X-Requested-With: XMLHttpRequest
Content-Length: 9122
Origin: https://banner9-registration.kfupm.edu.sa
Referer: https://banner9-registration.kfupm.edu.sa/StudentRegistrationSsb/ssb/classRegistration/classRegistration
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: same-origin
Priority: u=0
Te: trailers
Connection: keep-alive

{"create":[],"update":[{"addAuthorizationCrnMessage":null,"addAuthorizationCrnStatus":"INCOMPLETE","addDate":"01/06/2026","approvalOverride":"N","approvalReceivedIndicator":null,"approvalReceivedIndicatorHold":null,"attached":true,"attemptedHours":0,"authorizationCode":null,"billHour":3,"billHourInitial":3,"billHours":{"class":"net.hedtech.banner.student.registration.RegistrationCreditHour","creditHourHigh":null,"creditHourIndicator":null,"creditHourList":null,"creditHourLow":null},"block":null,"blockPermitOverride":null,"blockRuleSequenceNumber":null,"campus":"C","campusOverride":"N","capcOverride":"N","censusEnrollmentDate":"12/04/2026","class":"net.hedtech.banner.student.registration.RegistrationTemporaryView","classOverride":"N","cohortOverride":"N","collegeOverride":"N","completionDate":"05/10/2026","corqOverride":"N","courseContinuingEducationIndicator":"N","courseNumber":"208","courseReferenceNumber":"22729","courseRegistrationStatus":"RW","courseRegistrationStatusDescription":"Web Registered","courseTitle":"Differential Equations and Linear Algebra","creditHour":3,"creditHourInitial":3,"creditHours":{"class":"net.hedtech.banner.student.registration.RegistrationCreditHour","creditHourHigh":null,"creditHourIndicator":null,"creditHourList":null,"creditHourLow":null},"crnErrors":[],"dataOrigin":"Banner","degreeOverride":"N","departmentOverride":"N","dirty":false,"dirtyPropertyNames":[],"duplicateOverride":"N","durationUnit":null,"errorFlag":null,"errorLink":null,"errors":{"errors":[]},"grade":null,"gradeComment":null,"gradeDate":null,"gradeMid":null,"gradingMode":"S","gradingModeDescription":"Standard Grade Mode","id":25491123,"instructionalMethodDescription":null,"lastModified":"01/06/2026","level":"UG","levelDescription":"Undergraduate","levelOverride":"N","linkOverride":"N","majorOverride":"N","maxEnrollment":30,"message":null,"messageType":null,"mexcOverride":"N","newBlock":null,"newBlockRuleSequenceNumber":null,"numberOfUnits":null,"originalCourseRegistrationStatus":null,"originalRecordStatus":"N","originalVoiceResponseStatusType":null,"overrideDurationIndicator":false,"partOfTerm":"1","partOfTermDescription":"Full Term","permitOverrideUpdate":null,"preqOverride":"N","programOverride":"N","properties":{"billHourInitial":3,"registrationToDate":null,"errorFlag":null,"reservedKey":null,"courseContinuingEducationIndicator":"N","messageType":null,"selectedLevel":{"class":"net.hedtech.banner.student.registration.RegistrationLevel","description":null,"level":null},"courseRegistrationStatus":"RW","studentAttributeOverride":"N","approvalOverride":"N","classOverride":"N","subjectDescription":"Mathematics","campusOverride":"N","gradingMode":"S","selectedCreditHour":null,"block":null,"gradeComment":null,"sequenceNumber":"04","creditHourInitial":3,"courseRegistrationStatusDescription":"Web Registered","tuitionWaiverIndicator":"N","repeatOverride":"N","preqOverride":"N","creditHour":3,"majorOverride":"N","courseReferenceNumber":"22729","creditHours":{"class":"net.hedtech.banner.student.registration.RegistrationCreditHour","creditHourHigh":null,"creditHourIndicator":null,"creditHourList":null,"creditHourLow":null},"structuredRegistrationHeaderSequence":null,"scheduleType":"LEC","approvalReceivedIndicator":null,"grade":null,"programOverride":"N","startDate":"01/11/2026","addAuthorizationCrnMessage":null,"originalCourseRegistrationStatus":null,"registrationLevels":[],"authorizationCode":null,"attemptedHours":0,"selectedGradingMode":{"class":"net.hedtech.banner.student.registration.RegistrationGradingMode","description":null,"gradingMode":null},"waivHour":3,"levelDescription":"Undergraduate","selectedStartEndDate":{"class":"net.hedtech.banner.student.registration.RegistrationOlrStartEndDate","courseReferenceNumber":null,"durationUnit":null,"durationUnitDescription":null,"endDate":null,"numberOfUnits":null,"overrideDurationIndicator":null,"registrationDate":null,"sectionEndFromDate":null,"sectionEndToDate":null,"sectionStartFromDate":null,"sectionStartToDate":null,"startDate":null,"systemIn":null},"selectedOverride":null,"originalRecordStatus":"N","testOverride":"N","submitResultIndicator":null,"dataOrigin":"Banner","term":"202520","mexcOverride":"N","registrationOverrides":[],"levelOverride":"N","structuredRegistrationDetailSequence":null,"registrationFromDate":null,"waitCapacity":0,"blockRuleSequenceNumber":null,"sessionId":null,"billHours":{"class":"net.hedtech.banner.student.registration.RegistrationCreditHour","creditHourHigh":null,"creditHourIndicator":null,"creditHourList":null,"creditHourLow":null},"addDate":"01/06/2026","blockPermitOverride":null,"studyPathName":null,"addAuthorizationCrnStatus":"INCOMPLETE","studyPathKeySequence":null,"recordStatus":"N","numberOfUnits":null,"scheduleDescription":"Lecture","timeStatusHours":0,"selectedBillHour":null,"gradeDate":null,"registrationActions":[{"class":"net.hedtech.banner.student.registration.RegistrationAction","courseRegistrationStatus":"RW","description":"Web Registered","registrationStatusDate":null,"remove":false,"subActions":null,"voiceType":"R"},{"class":"net.hedtech.banner.student.registration.RegistrationAction","courseRegistrationStatus":"internal-remove","description":"Remove","registrationStatusDate":null,"remove":true,"subActions":null,"voiceType":null}],"errorLink":null,"maxEnrollment":30,"statusIndicator":"P","subject":"MATH","removeIndicator":null,"partOfTerm":"1","billHour":3,"sectionCourseTitle":"Differen. Eq & Linear Algebra","overrideDurationIndicator":false,"courseTitle":"Differential Equations and Linear Algebra","duplicateOverride":"N","durationUnit":null,"censusEnrollmentDate":"12/04/2026","corqOverride":"N","gradeMid":null,"level":"UG","instructionalMethodDescription":null,"campus":"C","degreeOverride":"N","waitOverride":"N","newBlock":null,"rpthOverride":"N","newBlockRuleSequenceNumber":null,"gradingModeDescription":"Standard Grade Mode","partOfTermDescription":"Full Term","lastModified":"01/06/2026","timeOverride":"N","registrationAuthorizationActiveCode":null,"linkOverride":"N","registrationStatusDate":"01/06/2026","departmentOverride":"N","crnErrors":[],"registrationGradingModes":[],"selectedStudyPath":{"class":"net.hedtech.banner.student.registration.RegistrationStudyPath","description":null,"keySequenceNumber":null},"approvalReceivedIndicatorHold":null,"capcOverride":"N","specialApproval":null,"courseNumber":"208","message":null,"cohortOverride":"N","statusDescription":"Pending","selectedAction":null,"collegeOverride":"N","originalVoiceResponseStatusType":null,"permitOverrideUpdate":null,"completionDate":"05/10/2026","voiceResponseStatusType":"R","registrationStudyPaths":[]},"recordStatus":"N","registrationActions":[{"class":"net.hedtech.banner.student.registration.RegistrationAction","courseRegistrationStatus":"RW","description":"Web Registered","registrationStatusDate":null,"remove":false,"subActions":null,"voiceType":"R"},{"class":"net.hedtech.banner.student.registration.RegistrationAction","courseRegistrationStatus":"internal-remove","description":"Remove","registrationStatusDate":null,"remove":true,"subActions":null,"voiceType":null}],"registrationAuthorizationActiveCode":null,"registrationFromDate":null,"registrationGradingModes":[],"registrationLevels":[],"registrationOverrides":[],"registrationStatusDate":"01/06/2026","registrationStudyPaths":[],"registrationToDate":null,"removeIndicator":null,"repeatOverride":"N","reservedKey":null,"rpthOverride":"N","scheduleDescription":"Lecture","scheduleType":"LEC","sectionCourseTitle":"Differen. Eq &amp; Linear Algebra","selectedAction":null,"selectedBillHour":null,"selectedCreditHour":null,"selectedGradingMode":{"class":"net.hedtech.banner.student.registration.RegistrationGradingMode","description":null,"gradingMode":null},"selectedLevel":{"class":"net.hedtech.banner.student.registration.RegistrationLevel","description":null,"level":null},"selectedOverride":null,"selectedStartEndDate":{"class":"net.hedtech.banner.student.registration.RegistrationOlrStartEndDate","courseReferenceNumber":null,"durationUnit":null,"durationUnitDescription":null,"endDate":null,"numberOfUnits":null,"overrideDurationIndicator":null,"registrationDate":null,"sectionEndFromDate":null,"sectionEndToDate":null,"sectionStartFromDate":null,"sectionStartToDate":null,"startDate":null,"systemIn":null},"selectedStudyPath":{"class":"net.hedtech.banner.student.registration.RegistrationStudyPath","description":null,"keySequenceNumber":null},"sequenceNumber":"04","sessionId":null,"specialApproval":null,"startDate":"01/11/2026","statusDescription":"Pending","statusIndicator":"P","structuredRegistrationDetailSequence":null,"structuredRegistrationHeaderSequence":null,"studentAttributeOverride":"N","studyPathKeySequence":null,"studyPathName":null,"subject":"MATH","subjectDescription":"Mathematics","submitResultIndicator":null,"term":"202520","testOverride":"N","timeOverride":"N","timeStatusHours":0,"tuitionWaiverIndicator":"N","version":0,"voiceResponseStatusType":"R","waitCapacity":0,"waitOverride":"N","waivHour":3}],"destroy":[],"uniqueSessionId":"yy10o1767691781029"}`

	req2, err := parseRawRequest(RawRequest2)
	if err != nil {
		return false, fmt.Errorf("failed to parse Req2: %v", err)
	}

	client2 := &http.Client{}
	resp2, err := client2.Do(req2)
	if err != nil {
		return false, fmt.Errorf("req2 failed: %v", err)
	}
	defer resp2.Body.Close()

	bodyBytes, err := io.ReadAll(resp2.Body)
	if err != nil {
		return false, err
	}

	bodyString := string(bodyBytes)

	if strings.Contains(bodyString, "Closed Section") {
		submit()
		fmt.Println("Response Preview:", bodyString)
		return false, nil
	}

	submit()
	fmt.Println("Response Preview:", bodyString)
	return true, nil
}

func submit() error {
	RawRequest1 := `POST /StudentRegistrationSsb/ssb/classRegistration/submitRegistration/batch HTTP/1.1
Host: banner9-registration.kfupm.edu.sa
Cookie: JSESSIONID=760032F45B31FD0370D2B9FAE3923EEE; _ga=GA1.3.1500780047.1767691760; _gid=GA1.3.2012430359.1767691760; _ga_EM00LC24FH=GS2.3.s1767691761$o1$g1$t1767691772$j49$l0$h0; KFUPM_Cookie=!F+vriIDI35PUtueAG8WAw/NalnvuQcHZoIEJ5QmvFhbsAFC9keeTENAtq/2tKv/e0XwdyTYNYGr9yLA=
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:146.0) Gecko/20100101 Firefox/146.0
Accept: application/json, text/javascript, */*; q=0.01
Accept-Language: en-US,en;q=0.5
Accept-Encoding: gzip, deflate, br
X-Synchronizer-Token: 3a81e06f-805f-47b9-a490-21156f23400c
Content-Type: application/json
X-Requested-With: XMLHttpRequest
Content-Length: 9247
Origin: https://banner9-registration.kfupm.edu.sa
Referer: https://banner9-registration.kfupm.edu.sa/StudentRegistrationSsb/ssb/classRegistration/classRegistration
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: same-origin
Priority: u=0
Te: trailers
Connection: keep-alive

{"create":[],"update":[{"addAuthorizationCrnMessage":null,"addAuthorizationCrnStatus":null,"addDate":"01/06/2026","approvalOverride":"N","approvalReceivedIndicator":null,"approvalReceivedIndicatorHold":null,"attached":true,"attemptedHours":0,"authorizationCode":null,"billHour":3,"billHourInitial":3,"billHours":{"class":"net.hedtech.banner.student.registration.RegistrationCreditHour","creditHourHigh":null,"creditHourIndicator":null,"creditHourList":null,"creditHourLow":null},"block":null,"blockPermitOverride":null,"blockRuleSequenceNumber":null,"campus":"C","campusOverride":"N","capcOverride":"N","censusEnrollmentDate":null,"class":"net.hedtech.banner.student.registration.RegistrationTemporaryView","classOverride":"N","cohortOverride":"N","collegeOverride":"N","completionDate":"05/10/2026","corqOverride":"N","courseContinuingEducationIndicator":"N","courseNumber":"208","courseReferenceNumber":"22729","courseRegistrationStatus":"RW","courseRegistrationStatusDescription":"Web Registered","courseTitle":"Differential Equations and Linear Algebra","creditHour":3,"creditHourInitial":3,"creditHours":{"class":"net.hedtech.banner.student.registration.RegistrationCreditHour","creditHourHigh":null,"creditHourIndicator":null,"creditHourList":null,"creditHourLow":null},"crnErrors":[{"class":"net.hedtech.banner.student.registration.RegistrationMessage","errorFlag":"F","message":"Closed Section ","messageType":"CLOS"},{"class":"net.hedtech.banner.student.registration.RegistrationMessage","errorFlag":"F","message":"Duplicate Course with Section 22486","messageType":"DUPL"}],"dataOrigin":"Banner","degreeOverride":"N","departmentOverride":"N","dirty":false,"dirtyPropertyNames":[],"duplicateOverride":"N","durationUnit":null,"errorFlag":"F","errorLink":null,"errors":{"errors":[]},"grade":null,"gradeComment":null,"gradeDate":null,"gradeMid":null,"gradingMode":"S","gradingModeDescription":"Standard Grade Mode","id":25491123,"instructionalMethodDescription":null,"lastModified":"01/06/2026","level":"UG","levelDescription":"Undergraduate","levelOverride":"N","linkOverride":"N","majorOverride":"N","maxEnrollment":null,"message":"Closed Section ","messageType":"CLOS","mexcOverride":"N","newBlock":null,"newBlockRuleSequenceNumber":null,"numberOfUnits":null,"originalCourseRegistrationStatus":null,"originalRecordStatus":"N","originalVoiceResponseStatusType":null,"overrideDurationIndicator":false,"partOfTerm":"1","partOfTermDescription":"Full Term","permitOverrideUpdate":null,"preqOverride":"N","programOverride":"N","properties":{"billHourInitial":3,"registrationToDate":null,"errorFlag":null,"reservedKey":null,"courseContinuingEducationIndicator":"N","messageType":null,"selectedLevel":{"class":"net.hedtech.banner.student.registration.RegistrationLevel","description":null,"level":null},"courseRegistrationStatus":"RW","studentAttributeOverride":"N","approvalOverride":"N","classOverride":"N","subjectDescription":"Mathematics","campusOverride":"N","gradingMode":"S","selectedCreditHour":null,"block":null,"gradeComment":null,"sequenceNumber":"04","creditHourInitial":3,"courseRegistrationStatusDescription":"Web Registered","tuitionWaiverIndicator":"N","repeatOverride":"N","preqOverride":"N","creditHour":3,"majorOverride":"N","courseReferenceNumber":"22729","creditHours":{"class":"net.hedtech.banner.student.registration.RegistrationCreditHour","creditHourHigh":null,"creditHourIndicator":null,"creditHourList":null,"creditHourLow":null},"structuredRegistrationHeaderSequence":null,"scheduleType":"LEC","approvalReceivedIndicator":null,"grade":null,"programOverride":"N","startDate":"01/11/2026","addAuthorizationCrnMessage":null,"originalCourseRegistrationStatus":null,"registrationLevels":[],"authorizationCode":null,"attemptedHours":0,"selectedGradingMode":{"class":"net.hedtech.banner.student.registration.RegistrationGradingMode","description":null,"gradingMode":null},"waivHour":3,"levelDescription":"Undergraduate","selectedStartEndDate":{"class":"net.hedtech.banner.student.registration.RegistrationOlrStartEndDate","courseReferenceNumber":null,"durationUnit":null,"durationUnitDescription":null,"endDate":null,"numberOfUnits":null,"overrideDurationIndicator":null,"registrationDate":null,"sectionEndFromDate":null,"sectionEndToDate":null,"sectionStartFromDate":null,"sectionStartToDate":null,"startDate":null,"systemIn":null},"selectedOverride":null,"originalRecordStatus":"N","testOverride":"N","submitResultIndicator":null,"dataOrigin":"Banner","term":"202520","mexcOverride":"N","registrationOverrides":[],"levelOverride":"N","structuredRegistrationDetailSequence":null,"registrationFromDate":null,"waitCapacity":0,"blockRuleSequenceNumber":null,"sessionId":null,"billHours":{"class":"net.hedtech.banner.student.registration.RegistrationCreditHour","creditHourHigh":null,"creditHourIndicator":null,"creditHourList":null,"creditHourLow":null},"addDate":"01/06/2026","blockPermitOverride":null,"studyPathName":null,"addAuthorizationCrnStatus":"INCOMPLETE","studyPathKeySequence":null,"recordStatus":"N","numberOfUnits":null,"scheduleDescription":"Lecture","timeStatusHours":0,"selectedBillHour":null,"gradeDate":null,"registrationActions":[{"class":"net.hedtech.banner.student.registration.RegistrationAction","courseRegistrationStatus":"RW","description":"Web Registered","registrationStatusDate":null,"remove":false,"subActions":null,"voiceType":"R"},{"class":"net.hedtech.banner.student.registration.RegistrationAction","courseRegistrationStatus":"internal-remove","description":"Remove","registrationStatusDate":null,"remove":true,"subActions":null,"voiceType":null}],"errorLink":null,"maxEnrollment":30,"statusIndicator":"P","subject":"MATH","removeIndicator":null,"partOfTerm":"1","billHour":3,"sectionCourseTitle":"Differen. Eq & Linear Algebra","overrideDurationIndicator":false,"courseTitle":"Differential Equations and Linear Algebra","duplicateOverride":"N","durationUnit":null,"censusEnrollmentDate":"12/04/2026","corqOverride":"N","gradeMid":null,"level":"UG","instructionalMethodDescription":null,"campus":"C","degreeOverride":"N","waitOverride":"N","newBlock":null,"rpthOverride":"N","newBlockRuleSequenceNumber":null,"gradingModeDescription":"Standard Grade Mode","partOfTermDescription":"Full Term","lastModified":"01/06/2026","timeOverride":"N","registrationAuthorizationActiveCode":null,"linkOverride":"N","registrationStatusDate":"01/06/2026","departmentOverride":"N","crnErrors":[],"registrationGradingModes":[],"selectedStudyPath":{"class":"net.hedtech.banner.student.registration.RegistrationStudyPath","description":null,"keySequenceNumber":null},"approvalReceivedIndicatorHold":null,"capcOverride":"N","specialApproval":null,"courseNumber":"208","message":null,"cohortOverride":"N","statusDescription":"Pending","selectedAction":null,"collegeOverride":"N","originalVoiceResponseStatusType":null,"permitOverrideUpdate":null,"completionDate":"05/10/2026","voiceResponseStatusType":"R","registrationStudyPaths":[]},"recordStatus":"N","registrationActions":[{"class":"net.hedtech.banner.student.registration.RegistrationAction","courseRegistrationStatus":"internal-remove","description":"Remove","registrationStatusDate":null,"remove":true,"subActions":null,"voiceType":null},{"class":"net.hedtech.banner.student.registration.RegistrationAction","courseRegistrationStatus":"RW","description":"Web Registered","registrationStatusDate":null,"remove":false,"subActions":null,"voiceType":"R"}],"registrationAuthorizationActiveCode":null,"registrationFromDate":null,"registrationGradingModes":[],"registrationLevels":[],"registrationOverrides":[],"registrationStatusDate":"01/06/2026","registrationStudyPaths":[],"registrationToDate":null,"removeIndicator":"Y","repeatOverride":"N","reservedKey":null,"rpthOverride":"N","scheduleDescription":"Lecture","scheduleType":"LEC","sectionCourseTitle":"Differen. Eq & Linear Algebra","selectedAction":{"class":"net.hedtech.banner.student.registration.RegistrationAction","courseRegistrationStatus":null,"description":null,"registrationStatusDate":null,"remove":false,"subActions":null,"voiceType":null},"selectedBillHour":null,"selectedCreditHour":null,"selectedGradingMode":{"class":"net.hedtech.banner.student.registration.RegistrationGradingMode","description":null,"gradingMode":null},"selectedLevel":{"class":"net.hedtech.banner.student.registration.RegistrationLevel","description":null,"level":null},"selectedOverride":null,"selectedStartEndDate":null,"selectedStudyPath":{"class":"net.hedtech.banner.student.registration.RegistrationStudyPath","description":null,"keySequenceNumber":null},"sequenceNumber":"04","sessionId":null,"specialApproval":null,"startDate":"01/11/2026","statusDescription":"Errors Preventing Registration","statusIndicator":"F","structuredRegistrationDetailSequence":null,"structuredRegistrationHeaderSequence":null,"studentAttributeOverride":"N","studyPathKeySequence":null,"studyPathName":null,"subject":"MATH","subjectDescription":"Mathematics","submitResultIndicator":"F","term":"202520","testOverride":"N","timeOverride":"N","timeStatusHours":0,"tuitionWaiverIndicator":"N","version":0,"voiceResponseStatusType":"R","waitCapacity":null,"waitOverride":"N","waivHour":3}],"destroy":[],"uniqueSessionId":"yy10o1767691781029"}`

	req1, err := parseRawRequest(RawRequest1)
	if err != nil {
		return err
	}

	// EXECUTE REQUEST 1
	client := &http.Client{}
	resp, err := client.Do(req1)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Helper to safely parse raw string, fix scheme, and recalc content-length
func parseRawRequest(raw string) (*http.Request, error) {
	// 1. Normalize Newlines
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	raw = strings.ReplaceAll(raw, "\n", "\r\n")

	// 2. Parse
	reader := bufio.NewReader(strings.NewReader(raw))
	req, err := http.ReadRequest(reader)
	if err != nil {
		return nil, err
	}

	// 3. Fix Protocol (CRITICAL STEP)
	req.RequestURI = ""
	req.URL.Scheme = "https" // Must be HTTPS for Banner9
	req.URL.Host = req.Host

	// 4. Recalculate Content-Length
	// ReadRequest uses the Content-Length header from the string,
	// which might be wrong after newline conversion.
	bodyBytes, _ := io.ReadAll(req.Body)
	req.ContentLength = int64(len(bodyBytes))
	req.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	// 5. Cleanup Headers
	req.Header.Del("Accept-Encoding") // Avoid binary response

	return req, nil
}

func sendDesktopAlert() {
	fmt.Println(">> TRIGGERING ALERTS...")

	// 1. Play Sound (Keep this, it works)
	beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)

	// 2. Try the Beeep notification (might fail/be silent)
	beeep.Alert("Course Available!", "Go register immediately!", "")

	// 3. FORCE A WINDOWS POPUP (The "Hard" Alert)
	// This runs a PowerShell command to open a message box in the middle of the screen.
	// It works even if Focus Assist is on.
	cmd := exec.Command("powershell", "-Command",
		"Add-Type -AssemblyName PresentationFramework;[System.Windows.MessageBox]::Show('THE COURSE IS OPEN! GO CHECK IT!', 'URGENT ALERT')")

	// We run this in a separate goroutine so it doesn't block the rest of the program
	// (Though for an alert, blocking is usually fine)
	_ = cmd.Start()
}

func sendEmail() {
	// Configuration
	from := "noureldinshaban4@gmail.com"
	password := "lsna gxef enxh jdgn" // Paste the App Password here
	to := "s202469600@kfupm.edu.sa"   // You can send it to yourself
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message
	message := []byte("Subject: Reaminig Seats Found!\r\n" +
		"\r\n" +
		"الكورس فتح يا بيروووو سجل بسرعة قبل ما يقفل تاني!\r\n")

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return
	}
	fmt.Println("Email notification sent successfully.")
}

// Helper function to Check availability of course

// Helper function to drop another course

// func removeCourse() (bool, error) {

// 	RawRequest := `POST /StudentRegistrationSsb/ssb/classRegistration/submitRegistration/batch HTTP/1.1
// Host: banner9-registration.kfupm.edu.sa
// Cookie: JSESSIONID=1FDB4A0939A9BF73DBC4941B5E0786E8; _ga=GA1.3.1500780047.1767691760; _gid=GA1.3.2012430359.1767691760; _ga_EM00LC24FH=GS2.3.s1767691761$o1$g1$t1767691772$j49$l0$h0; KFUPM_Cookie=!spuPm3HhYzemQW+AG8WAw/NalnvuQeZEArkok2HrIYrSQSNY+Vwgukj4wa7IXoTiBlbbSKciFVlQ0Ek=
// User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:146.0) Gecko/20100101 Firefox/146.0
// Accept: application/json, text/javascript, */*; q=0.01
// Accept-Language: en-US,en;q=0.5
// Accept-Encoding: gzip, deflate, br
// X-Synchronizer-Token: 1ba389a6-081e-4bec-822e-f5dae6348b47
// Content-Type: application/json
// X-Requested-With: XMLHttpRequest
// Content-Length: 9090
// Origin: https://banner9-registration.kfupm.edu.sa
// Referer: https://banner9-registration.kfupm.edu.sa/StudentRegistrationSsb/ssb/classRegistration/classRegistration
// Sec-Fetch-Dest: empty
// Sec-Fetch-Mode: cors
// Sec-Fetch-Site: same-origin
// Priority: u=0
// Te: trailers
// Connection: keep-alive

// {"create":[],"update":[{"addAuthorizationCrnMessage":null,"addAuthorizationCrnStatus":null,"addDate":"01/06/2026","approvalOverride":"N","approvalReceivedIndicator":null,"approvalReceivedIndicatorHold":null,"attached":true,"attemptedHours":0,"authorizationCode":null,"billHour":3,"billHourInitial":3,"billHours":{"class":"net.hedtech.banner.student.registration.RegistrationCreditHour","creditHourHigh":null,"creditHourIndicator":null,"creditHourList":null,"creditHourLow":null},"block":null,"blockPermitOverride":null,"blockRuleSequenceNumber":null,"campus":"C","campusOverride":"N","capcOverride":"N","censusEnrollmentDate":null,"class":"net.hedtech.banner.student.registration.RegistrationTemporaryView","classOverride":"N","cohortOverride":"N","collegeOverride":"N","completionDate":"05/10/2026","corqOverride":"N","courseContinuingEducationIndicator":"N","courseNumber":"208","courseReferenceNumber":"22729","courseRegistrationStatus":"RW","courseRegistrationStatusDescription":"Web Registered","courseTitle":"Differential Equations and Linear Algebra","creditHour":3,"creditHourInitial":3,"creditHours":{"class":"net.hedtech.banner.student.registration.RegistrationCreditHour","creditHourHigh":null,"creditHourIndicator":null,"creditHourList":null,"creditHourLow":null},"crnErrors":[{"class":"net.hedtech.banner.student.registration.RegistrationMessage","errorFlag":"F","message":"Closed Section ","messageType":"CLOS"}],"dataOrigin":"Banner","degreeOverride":"N","departmentOverride":"N","dirty":false,"dirtyPropertyNames":[],"duplicateOverride":"N","durationUnit":null,"errorFlag":"F","errorLink":null,"errors":{"errors":[]},"grade":null,"gradeComment":null,"gradeDate":null,"gradeMid":null,"gradingMode":"S","gradingModeDescription":"Standard Grade Mode","id":25478533,"instructionalMethodDescription":null,"lastModified":"01/06/2026","level":"UG","levelDescription":"Undergraduate","levelOverride":"N","linkOverride":"N","majorOverride":"N","maxEnrollment":null,"message":"Closed Section ","messageType":"CLOS","mexcOverride":"N","newBlock":null,"newBlockRuleSequenceNumber":null,"numberOfUnits":null,"originalCourseRegistrationStatus":null,"originalRecordStatus":"N","originalVoiceResponseStatusType":null,"overrideDurationIndicator":false,"partOfTerm":"1","partOfTermDescription":"Full Term","permitOverrideUpdate":null,"preqOverride":"N","programOverride":"N","properties":{"billHourInitial":3,"registrationToDate":null,"errorFlag":null,"reservedKey":null,"courseContinuingEducationIndicator":"N","messageType":null,"courseRegistrationStatus":"RW","studentAttributeOverride":"N","approvalOverride":"N","classOverride":"N","selectedLevel":{"class":"net.hedtech.banner.student.registration.RegistrationLevel","description":null,"level":null},"subjectDescription":"Mathematics","campusOverride":"N","gradingMode":"S","selectedCreditHour":null,"block":null,"gradeComment":null,"sequenceNumber":"04","creditHourInitial":3,"courseRegistrationStatusDescription":"Web Registered","tuitionWaiverIndicator":"N","repeatOverride":"N","preqOverride":"N","creditHour":3,"majorOverride":"N","courseReferenceNumber":"22729","structuredRegistrationHeaderSequence":null,"creditHours":{"class":"net.hedtech.banner.student.registration.RegistrationCreditHour","creditHourHigh":null,"creditHourIndicator":null,"creditHourList":null,"creditHourLow":null},"scheduleType":"LEC","approvalReceivedIndicator":null,"grade":null,"programOverride":"N","startDate":"01/11/2026","addAuthorizationCrnMessage":null,"originalCourseRegistrationStatus":null,"registrationLevels":[],"authorizationCode":null,"attemptedHours":0,"selectedGradingMode":{"class":"net.hedtech.banner.student.registration.RegistrationGradingMode","description":null,"gradingMode":null},"waivHour":3,"levelDescription":"Undergraduate","selectedStartEndDate":{"class":"net.hedtech.banner.student.registration.RegistrationOlrStartEndDate","courseReferenceNumber":null,"durationUnit":null,"durationUnitDescription":null,"endDate":null,"numberOfUnits":null,"overrideDurationIndicator":null,"registrationDate":null,"sectionEndFromDate":null,"sectionEndToDate":null,"sectionStartFromDate":null,"sectionStartToDate":null,"startDate":null,"systemIn":null},"selectedOverride":null,"originalRecordStatus":"N","testOverride":"N","submitResultIndicator":null,"dataOrigin":"Banner","term":"202520","mexcOverride":"N","registrationOverrides":[],"levelOverride":"N","structuredRegistrationDetailSequence":null,"registrationFromDate":null,"waitCapacity":0,"blockRuleSequenceNumber":null,"sessionId":null,"addDate":"01/06/2026","billHours":{"class":"net.hedtech.banner.student.registration.RegistrationCreditHour","creditHourHigh":null,"creditHourIndicator":null,"creditHourList":null,"creditHourLow":null},"blockPermitOverride":null,"studyPathName":null,"addAuthorizationCrnStatus":"INCOMPLETE","studyPathKeySequence":null,"recordStatus":"N","numberOfUnits":null,"scheduleDescription":"Lecture","timeStatusHours":0,"selectedBillHour":null,"gradeDate":null,"registrationActions":[{"class":"net.hedtech.banner.student.registration.RegistrationAction","courseRegistrationStatus":"RW","description":"Web Registered","registrationStatusDate":null,"remove":false,"subActions":null,"voiceType":"R"},{"class":"net.hedtech.banner.student.registration.RegistrationAction","courseRegistrationStatus":"internal-remove","description":"Remove","registrationStatusDate":null,"remove":true,"subActions":null,"voiceType":null}],"errorLink":null,"maxEnrollment":30,"statusIndicator":"P","subject":"MATH","removeIndicator":null,"partOfTerm":"1","billHour":3,"sectionCourseTitle":"Differen. Eq & Linear Algebra","overrideDurationIndicator":false,"courseTitle":"Differential Equations and Linear Algebra","duplicateOverride":"N","durationUnit":null,"censusEnrollmentDate":"12/04/2026","corqOverride":"N","gradeMid":null,"level":"UG","instructionalMethodDescription":null,"campus":"C","degreeOverride":"N","waitOverride":"N","newBlock":null,"rpthOverride":"N","newBlockRuleSequenceNumber":null,"gradingModeDescription":"Standard Grade Mode","partOfTermDescription":"Full Term","lastModified":"01/06/2026","timeOverride":"N","registrationAuthorizationActiveCode":null,"linkOverride":"N","registrationStatusDate":"01/06/2026","departmentOverride":"N","crnErrors":[],"selectedStudyPath":{"class":"net.hedtech.banner.student.registration.RegistrationStudyPath","description":null,"keySequenceNumber":null},"registrationGradingModes":[],"approvalReceivedIndicatorHold":null,"capcOverride":"N","specialApproval":null,"courseNumber":"208","message":null,"cohortOverride":"N","statusDescription":"Pending","selectedAction":null,"collegeOverride":"N","originalVoiceResponseStatusType":null,"permitOverrideUpdate":null,"completionDate":"05/10/2026","voiceResponseStatusType":"R","registrationStudyPaths":[]},"recordStatus":"N","registrationActions":[{"class":"net.hedtech.banner.student.registration.RegistrationAction","courseRegistrationStatus":"internal-remove","description":"Remove","registrationStatusDate":null,"remove":true,"subActions":null,"voiceType":null},{"class":"net.hedtech.banner.student.registration.RegistrationAction","courseRegistrationStatus":"RW","description":"Web Registered","registrationStatusDate":null,"remove":false,"subActions":null,"voiceType":"R"}],"registrationAuthorizationActiveCode":null,"registrationFromDate":null,"registrationGradingModes":[],"registrationLevels":[],"registrationOverrides":[],"registrationStatusDate":"01/06/2026","registrationStudyPaths":[],"registrationToDate":null,"removeIndicator":"Y","repeatOverride":"N","reservedKey":null,"rpthOverride":"N","scheduleDescription":"Lecture","scheduleType":"LEC","sectionCourseTitle":"Differen. Eq & Linear Algebra","selectedAction":{"class":"net.hedtech.banner.student.registration.RegistrationAction","courseRegistrationStatus":null,"description":null,"registrationStatusDate":null,"remove":false,"subActions":null,"voiceType":null},"selectedBillHour":null,"selectedCreditHour":null,"selectedGradingMode":{"class":"net.hedtech.banner.student.registration.RegistrationGradingMode","description":null,"gradingMode":null},"selectedLevel":{"class":"net.hedtech.banner.student.registration.RegistrationLevel","description":null,"level":null},"selectedOverride":null,"selectedStartEndDate":null,"selectedStudyPath":{"class":"net.hedtech.banner.student.registration.RegistrationStudyPath","description":null,"keySequenceNumber":null},"sequenceNumber":"04","sessionId":null,"specialApproval":null,"startDate":"01/11/2026","statusDescription":"Errors Preventing Registration","statusIndicator":"F","structuredRegistrationDetailSequence":null,"structuredRegistrationHeaderSequence":null,"studentAttributeOverride":"N","studyPathKeySequence":null,"studyPathName":null,"subject":"MATH","subjectDescription":"Mathematics","submitResultIndicator":"F","term":"202520","testOverride":"N","timeOverride":"N","timeStatusHours":0,"tuitionWaiverIndicator":"N","version":0,"voiceResponseStatusType":"R","waitCapacity":null,"waitOverride":"N","waivHour":3}],"destroy":[],"uniqueSessionId":"yy10o1767691781029"}`

// 	rawReq := strings.ReplaceAll(RawRequest, "\n", "\r\n")

// 	reader := bufio.NewReader(strings.NewReader(rawReq))

// 	req, err := http.ReadRequest(reader)
// 	if err != nil {
// 		return false, err
// 	}

// 	req.RequestURI = ""
// 	req.URL.Scheme = "http"
// 	req.URL.Host = req.Host

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return false, err
// 	}
// 	defer resp.Body.Close()

// 	return true, nil
// }

// // Helper function to verfiy the modification of schedule

// func verify() (bool, error) {
// 	RawRequest := `GET /bwskfreg.P_ViewSchedule HTTP/1.1
// Host: banner.university.edu
// User-Agent: Mozilla/5.0
// Cookie: SESSID=xyz123;
// Accept: text/html`

// 	rawReq := strings.ReplaceAll(RawRequest, "\n", "\r\n")

// 	reader := bufio.NewReader(strings.NewReader(rawReq))

// 	req, err := http.ReadRequest(reader)
// 	if err != nil {
// 		return false, err
// 	}
// 	req.RequestURI = ""
// 	req.URL.Scheme = "http"
// 	req.URL.Host = req.Host
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return false, err
// 	}
// 	defer resp.Body.Close()

// 	bodyBytes, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return false, err
// 	}

// 	bodyString := string(bodyBytes)

// 	if resp.StatusCode == http.StatusOK && strings.Contains(bodyString, "Current Schedule") {
// 		return true, nil
// 	}
// 	return false, nil
// }
