package def

var (
	AllValueTables = map[string]map[Range]string{
		"MDFieldFormat":                                        MDFieldFormat,
		"MinQueueSizeLogFieldValues":                           MinQueueSizeLogFieldValues,
		"ReceiveDelayFieldValues":                              ReceiveDelayFieldValues,
		"PollTimeoutFieldValues":                               PollTimeoutFieldValues,
		"NumElementsValueDefinitions":                          NumElementsValueDefinitions,
		"ReceiveWindowValueDefinitions":                        ReceiveWindowValueDefinitions,
		"InitTTLValueDefinitions":                              InitTTLValueDefinitions,
		"StepResolutionValues":                                 StepResolutionValues,
		"NumberOfStepsValues":                                  NumberOfStepsValues,
		"PublishFriendshipCredentialFlagValues":                PublishFriendshipCredentialFlagValues,
		"PublishTTLValues":                                     PublishTTLValues,
		"DefaultTTLValues":                                     DefaultTTLValues,
		"RelayValues":                                          RelayValues,
		"AttentionTimerValues":                                 AttentionTimerValues,
		"SecureNetworkBeaconValues":                            SecureNetworkBeaconValues,
		"GATTProxyValues":                                      GATTProxyValues,
		"NodeIdentityValues":                                   NodeIdentityValues,
		"FriendValues":                                         FriendValues,
		"KeyRefreshPhaseStateValues":                           KeyRefreshPhaseStateValues,
		"TestIDValues":                                         TestIDValues,
		"FaultValues":                                          FaultValues,
		"HeartbeatPublicationCountLogValues":                   HeartbeatPublicationCountLogValues,
		"HeartbeatPublicationPeriodLogValues":                  HeartbeatPublicationPeriodLogValues,
		"HeartbeatPublicationTTLValues":                        HeartbeatPublicationTTLValues,
		"HeartbeatSubscriptionCountValues":                     HeartbeatSubscriptionCountValues,
		"HeartbeatSubscriptionPeriodValues":                    HeartbeatSubscriptionPeriodValues,
		"HeartbeatSubscriptionMinTTLValues":                    HeartbeatSubscriptionMinTTLValues,
		"HeartbeatSubscriptionMaxTTLValues":                    HeartbeatSubscriptionMaxTTLValues,
		"PollTimeoutTimerValues":                               PollTimeoutTimerValues,
		"SummaryOfStatusCodes":                                 SummaryOfStatusCodes,
		"GenericProvisioningControlFormatFieldValues":          GenericProvisioningControlFormatFieldValues,
		"NumberOfElementsFieldValues":                          NumberOfElementsFieldValues,
		"OutputOOBSizeFieldValues":                             OutputOOBSizeFieldValues,
		"InputOOBSizeFieldValues":                              InputOOBSizeFieldValues,
		"AlgorithmFieldValues":                                 AlgorithmFieldValues,
		"PublicKeyFieldValues":                                 PublicKeyFieldValues,
		"AuthenticationMethodFieldValues":                      AuthenticationMethodFieldValues,
		"OutputOOBActionValuesForTheAuthenticationActionField": OutputOOBActionValuesForTheAuthenticationActionField,
		"OutputOOBSizeValuesForTheAuthenticationSizeField":     OutputOOBSizeValuesForTheAuthenticationSizeField,
		"InputOOBActionValuesForTheAuthenticationActionField":  InputOOBActionValuesForTheAuthenticationActionField,
		"InputOOBSizeValuesForTheAuthenticationSizeField":      InputOOBSizeValuesForTheAuthenticationSizeField,
		"SARFieldValues":                                       SARFieldValues,
		"IdentificationTypeValues":                             IdentificationTypeValues,

		"GenericOnOffStates":                                              GenericOnOffStates,
		"GenericLevelStates":                                              GenericLevelStates,
		"DefaultTransitionStepResolutionValues":                           DefaultTransitionStepResolutionValues,
		"DefaultTransitionNumberOfStepsValues":                            DefaultTransitionNumberOfStepsValues,
		"GenericOnPowerUpStates":                                          GenericOnPowerUpStates,
		"GenericPowerActualStates":                                        GenericPowerActualStates,
		"GenericPowerLastStates":                                          GenericPowerLastStates,
		"GenericPowerDefaultStates":                                       GenericPowerDefaultStates,
		"GenericPowerMinAndGenericPowerMaxStates":                         GenericPowerMinAndGenericPowerMaxStates,
		"GenericBatteryLevelStates":                                       GenericBatteryLevelStates,
		"GenericBatteryTimeToDischargeStates":                             GenericBatteryTimeToDischargeStates,
		"GenericBatteryTimeToChargeStates":                                GenericBatteryTimeToChargeStates,
		"GenericBatteryFlagsPresenceStates":                               GenericBatteryFlagsPresenceStates,
		"GenericBatteryFlagsIndicatorStates":                              GenericBatteryFlagsIndicatorStates,
		"GenericBatteryFlagsChargingStates":                               GenericBatteryFlagsChargingStates,
		"GenericBatteryFlagsServiceabilityStates":                         GenericBatteryFlagsServiceabilityStates,
		"SensorPositiveToleranceStates":                                   SensorPositiveToleranceStates,
		"SensorNegativeToleranceStates":                                   SensorNegativeToleranceStates,
		"SensorSamplingFunctions":                                         SensorSamplingFunctions,
		"SensorDataFormatValues":                                          SensorDataFormatValues,
		"SceneNumberValues":                                               SceneNumberValues,
		"YearFieldValues":                                                 YearFieldValues,
		"DayFieldValues":                                                  DayFieldValues,
		"HourFieldValues":                                                 HourFieldValues,
		"MinuteFieldValues":                                               MinuteFieldValues,
		"SecondFieldValues":                                               SecondFieldValues,
		"ActionFieldValues":                                               ActionFieldValues,
		"SceneNumberFieldValues":                                          SceneNumberFieldValues,
		"StatusCodeValues":                                                StatusCodeValues,
		"LightLightnessLinearStates":                                      LightLightnessLinearStates,
		"LightLightnessActualStates":                                      LightLightnessActualStates,
		"LightLightnessLastStates":                                        LightLightnessLastStates,
		"LightLightnessDefaultStates":                                     LightLightnessDefaultStates,
		"LightLightnessMinAndLightLightnessMaxStates":                     LightLightnessMinAndLightLightnessMaxStates,
		"LightCTLTemperatureStates":                                       LightCTLTemperatureStates,
		"LightCTLTemperatureDefaultStates":                                LightCTLTemperatureDefaultStates,
		"LightCTLTemperatureRangeMinAndLightCTLTemperatureRangeMaxStates": LightCTLTemperatureRangeMinAndLightCTLTemperatureRangeMaxStates,
		"LightCTLDeltaUVStates":                                           LightCTLDeltaUVStates,
		"LightCTLDeltaUVDefaultStates":                                    LightCTLDeltaUVDefaultStates,
		"LightCTLLightnessStates":                                         LightCTLLightnessStates,
		"LightHSLHueStates":                                               LightHSLHueStates,
		"LightHSLHueDefaultStates":                                        LightHSLHueDefaultStates,
		"LightHSLHueMinAndLightHSLHueMaxStates":                           LightHSLHueMinAndLightHSLHueMaxStates,
		"LightHSLSaturationStates":                                        LightHSLSaturationStates,
		"LightHSLSaturationDefaultStates":                                 LightHSLSaturationDefaultStates,
		"LightHSLSaturationMinAndLightHSLSaturationMaxStates":             LightHSLSaturationMinAndLightHSLSaturationMaxStates,
		"LightHSLLightnessStates":                                         LightHSLLightnessStates,
		"LightXyLXStates":                                                 LightXyLXStates,
		"LightXyLXDefaultStates":                                          LightXyLXDefaultStates,
		"LightXyLXMinAndLightXyLXMaxStates":                               LightXyLXMinAndLightXyLXMaxStates,
		"LightXyLYStates":                                                 LightXyLYStates,
		"LightXyLYDefaultStates":                                          LightXyLYDefaultStates,
		"LightXyLYMinAndLightXyLYMaxStates":                               LightXyLYMinAndLightXyLYMaxStates,
		"LightXyLLightnessStates":                                         LightXyLLightnessStates,
		"LightLCModeStates":                                               LightLCModeStates,
		"LightLCOccupancyModeStates":                                      LightLCOccupancyModeStates,
		"LightLCLightOnOffStates":                                         LightLCLightOnOffStates,
		"LightLCOccupancyStates":                                          LightLCOccupancyStates,
		"LightLCAmbientLuxLevelStates":                                    LightLCAmbientLuxLevelStates,
		"LightLCLinearOutputStates":                                       LightLCLinearOutputStates,
		"LightLCAmbientLuxLevelOnStates":                                  LightLCAmbientLuxLevelOnStates,
		"LightLCAmbientLuxLevelProlongStates":                             LightLCAmbientLuxLevelProlongStates,
		"LightLCAmbientLuxLevelStandbyStates":                             LightLCAmbientLuxLevelStandbyStates,
	}
	// Table 3.20: MD field format
	MDFieldFormat = map[Range]string{
		Range{0, 0}:   "Friend Queue is empty",
		Range{1, 1}:   "Friend Queue is not empty",
		Range{2, 255}: "Prohibited",
	}

	// Table 3.25: MinQueueSizeLog field values
	MinQueueSizeLogFieldValues = map[Range]string{
		Range{0, 0}: "Prohibited",
		Range{1, 1}: "N = 2",
		Range{2, 2}: "N = 4",
		Range{3, 3}: "N = 8",
		Range{4, 4}: "N = 16",
		Range{5, 5}: "N = 32",
		Range{6, 6}: "N = 64",
		Range{7, 7}: "N = 128",
	}

	// Table 3.26: ReceiveDelay field values
	ReceiveDelayFieldValues = map[Range]string{
		Range{0, 9}:    "Prohibited",
		Range{10, 255}: "Receive Delay in units of 1 millisecond",
	}

	// Table 3.27: PollTimeout field values
	PollTimeoutFieldValues = map[Range]string{
		Range{0, 9}:              "Prohibited",
		Range{10, 3455999}:       "PollTimeout in units of 100 milliseconds",
		Range{3456000, 16777215}: "Prohibited",
	}

	// Table 3.28: NumElements value definitions
	NumElementsValueDefinitions = map[Range]string{
		Range{0, 0}:   "Prohibited",
		Range{1, 255}: "Number of elements",
	}

	// Table 3.30: ReceiveWindow value definitions
	ReceiveWindowValueDefinitions = map[Range]string{
		Range{0, 0}:   "Prohibited",
		Range{1, 255}: "Receive Window in units of 1 millisecond",
	}

	// Table 3.37: InitTTL value definitions
	InitTTLValueDefinitions = map[Range]string{
		Range{0, 127}: "Initial TTL when sending a message",
	}

	// Table 4.6: Step Resolution values
	StepResolutionValues = map[Range]string{
		Range{0, 0}: "The Step Resolution is 100 milliseconds",
		Range{1, 1}: "The Step Resolution is 1 second",
		Range{2, 2}: "The Step Resolution is 10 seconds",
		Range{3, 3}: "The Step Resolution is 10 minutes",
	}

	// Table 4.7: Number of Steps values
	NumberOfStepsValues = map[Range]string{
		Range{0, 0}:  "Publish Period is disabled",
		Range{1, 63}: "The number of steps",
	}

	// Table 4.8: Publish Friendship Credential Flag values
	PublishFriendshipCredentialFlagValues = map[Range]string{
		Range{0, 0}: "Master security material is used for Publishing",
		Range{1, 1}: "Friendship security material is used for Publishing",
	}

	// Table 4.9: Publish TTL values
	PublishTTLValues = map[Range]string{
		Range{0, 127}:   "The Publish TTL value, represented as a 1-octet integer",
		Range{128, 254}: "Prohibited",
		Range{255, 255}: "Use Default TTL",
	}

	// Table 4.10: Default TTL values
	DefaultTTLValues = map[Range]string{
		Range{0, 0}:     "The Default TTL state",
		Range{2, 127}:   "The Default TTL state",
		Range{1, 1}:     "Prohibited",
		Range{128, 255}: "Prohibited",
	}

	// Table 4.11: Relay values
	RelayValues = map[Range]string{
		Range{0, 0}:   "The node support Relay feature that is disabled",
		Range{1, 1}:   "The node supports Relay feature that is enabled",
		Range{2, 2}:   "Relay feature is not supported",
		Range{3, 255}: "Prohibited",
	}

	// Table 4.12: Attention Timer values
	AttentionTimerValues = map[Range]string{
		Range{0, 0}:   "Off",
		Range{1, 255}: "On, remaining time in seconds",
	}

	// Table 4.13: Secure Network Beacon values
	SecureNetworkBeaconValues = map[Range]string{
		Range{0, 0}:   "The node is not broadcasting a Secure Network beacon",
		Range{1, 1}:   "The node is broadcasting a Secure Network beacon",
		Range{2, 255}: "Prohibited",
	}

	// Table 4.14: GATT Proxy values
	GATTProxyValues = map[Range]string{
		Range{0, 0}:   "The Mesh Proxy Service is running, Proxy feature is disabled",
		Range{1, 1}:   "The Mesh Proxy Service is running, Proxy feature is enabled",
		Range{2, 2}:   "The Mesh Proxy Service is not supported, Proxy feature is not supported",
		Range{3, 255}: "Prohibited",
	}

	// Table 4.15: Node Identity values
	NodeIdentityValues = map[Range]string{
		Range{0, 0}:   "Node Identity for a subnet is stopped",
		Range{1, 1}:   "Node Identity for a subnet is running",
		Range{2, 2}:   "Node Identity is not supported",
		Range{3, 255}: "Prohibited",
	}

	// Table 4.16: Friend values
	FriendValues = map[Range]string{
		Range{0, 0}:   "The node supports Friend feature that is disabled",
		Range{1, 1}:   "The node supports Friend feature that is enabled",
		Range{2, 2}:   "The Friend feature is not supported",
		Range{3, 255}: "Prohibited",
	}

	// Table 4.17: Key Refresh Phase state values
	KeyRefreshPhaseStateValues = map[Range]string{
		Range{0, 0}:   "Normal operation; Key Refresh procedure is not active",
		Range{1, 1}:   "First phase of Key Refresh procedure",
		Range{2, 2}:   "Second phase of Key Refresh procedure",
		Range{3, 255}: "Prohibited",
	}

	// Table 4.20: Test ID values
	TestIDValues = map[Range]string{
		Range{0, 0}:   "Standard test",
		Range{1, 255}: "Vendor specific test",
	}

	// Table 4.21: Fault values
	FaultValues = map[Range]string{
		Range{10, 10}:   "No Load Error",
		Range{11, 11}:   "Overload Warning",
		Range{12, 12}:   "Overload Error",
		Range{13, 13}:   "Overheat Warning",
		Range{14, 14}:   "Overheat Error",
		Range{15, 15}:   "Condensation Warning",
		Range{16, 16}:   "Condensation Error",
		Range{17, 17}:   "Vibration Warning",
		Range{18, 18}:   "Vibration Error",
		Range{19, 19}:   "Configuration Warning",
		Range{20, 20}:   "Configuration Error",
		Range{21, 21}:   "Element Not Calibrated Warning",
		Range{22, 22}:   "Element Not Calibrated Error",
		Range{23, 23}:   "Memory Warning",
		Range{24, 24}:   "Memory Error",
		Range{25, 25}:   "Self-Test Warning",
		Range{26, 26}:   "Self-Test Error",
		Range{27, 27}:   "Input Too Low Warning",
		Range{28, 28}:   "Input Too Low Error",
		Range{29, 29}:   "Input Too High Warning",
		Range{30, 30}:   "Input Too High Error",
		Range{31, 31}:   "Input No Change Warning",
		Range{32, 32}:   "Input No Change Error",
		Range{33, 33}:   "Actuator Blocked Warning",
		Range{34, 34}:   "Actuator Blocked Error",
		Range{35, 35}:   "Housing Opened Warning",
		Range{36, 36}:   "Housing Opened Error",
		Range{37, 37}:   "Tamper Warning",
		Range{38, 38}:   "Tamper Error",
		Range{39, 39}:   "Device Moved Warning",
		Range{40, 40}:   "Device Moved Error",
		Range{41, 41}:   "Device Dropped Warning",
		Range{42, 42}:   "Device Dropped Error",
		Range{43, 43}:   "Overflow Warning",
		Range{44, 44}:   "Overflow Error",
		Range{45, 45}:   "Empty Warning",
		Range{46, 46}:   "Empty Error",
		Range{47, 47}:   "Internal Bus Warning",
		Range{48, 48}:   "Internal Bus Error",
		Range{49, 49}:   "Mechanism Jammed Warning",
		Range{50, 50}:   "Mechanism Jammed Error",
		Range{51, 127}:  "Reserved for Future Use",
		Range{128, 255}: "Vendor Specific Warning / Error",
	}

	// Table 4.24: Heartbeat Publication Count Log values
	HeartbeatPublicationCountLogValues = map[Range]string{
		Range{0, 0}:     "Heartbeat messages are not being sent periodically",
		Range{1, 17}:    "Number of Heartbeat messages, 2(n-1), that remain to be sent",
		Range{18, 254}:  "Prohibited",
		Range{255, 255}: "Heartbeat messages are being sent indefinitely",
	}

	// Table 4.25: Heartbeat Publication Period Log values
	HeartbeatPublicationPeriodLogValues = map[Range]string{
		Range{0, 0}:    "Heartbeat messages are not being sent periodically",
		Range{1, 17}:   "Period in 2(n-1) seconds for sending periodical Heartbeat messages",
		Range{18, 255}: "Prohibited",
	}

	// Table 4.26: Heartbeat Publication TTL values
	HeartbeatPublicationTTLValues = map[Range]string{
		Range{0, 127}:   "The Heartbeat Publication TTL state",
		Range{128, 255}: "Prohibited",
	}

	// Table 4.28: Heartbeat Subscription Count values
	HeartbeatSubscriptionCountValues = map[Range]string{
		Range{0, 65534}:     "Number of Heartbeat messages received",
		Range{65535, 65535}: "More than 0xFFFE messages have been received",
	}

	// Table 4.29: Heartbeat Subscription Period values
	HeartbeatSubscriptionPeriodValues = map[Range]string{
		Range{0, 0}:    "Heartbeat messages are not being processed",
		Range{1, 17}:   "Remaining period in 2(n-1) seconds for processing periodical Heartbeat messages",
		Range{18, 255}: "Prohibited",
	}

	// Table 4.30: Heartbeat Subscription Min TTL values
	HeartbeatSubscriptionMinTTLValues = map[Range]string{
		Range{0, 0}:     "No Heartbeat messages have been received",
		Range{1, 127}:   "The Heartbeat Subscription Min Hops state",
		Range{128, 255}: "Prohibited",
	}

	// Table 4.31: Heartbeat Subscription Max TTL values
	HeartbeatSubscriptionMaxTTLValues = map[Range]string{
		Range{0, 0}:     "No Heartbeat messages have been received.",
		Range{1, 127}:   "The Heartbeat Subscription Max Hops state",
		Range{128, 255}: "Prohibited",
	}

	// Table 4.32: PollTimeout Timer values
	PollTimeoutTimerValues = map[Range]string{
		Range{0, 0}:              "The node is no longer a Friend node of the Low Power node identified by the LPNAddress",
		Range{1, 9}:              "Prohibited",
		Range{10, 3455999}:       "The PollTimeout timer value in units of 100 milliseconds",
		Range{3456000, 16777215}: "Prohibited",
	}

	// Table 4.108: Summary of status codes
	SummaryOfStatusCodes = map[Range]string{
		Range{0, 0}:    "Success",
		Range{1, 1}:    "Invalid Address",
		Range{2, 2}:    "Invalid Model",
		Range{3, 3}:    "Invalid AppKey Index",
		Range{4, 4}:    "Invalid NetKey Index",
		Range{5, 5}:    "Insufficient Resources",
		Range{6, 6}:    "Key Index Already Stored",
		Range{7, 7}:    "Invalid Publish Parameters",
		Range{8, 8}:    "Not a Subscribe Model",
		Range{9, 9}:    "Storage Failure",
		Range{10, 10}:  "Feature Not Supported",
		Range{11, 11}:  "Cannot Update",
		Range{12, 12}:  "Cannot Remove",
		Range{13, 13}:  "Cannot Bind",
		Range{14, 14}:  "Temporarily Unable to Change State",
		Range{15, 15}:  "Cannot Set",
		Range{16, 16}:  "Unspecified Error",
		Range{17, 17}:  "Invalid Binding",
		Range{18, 255}: "RFU",
	}

	// Table 5.4: Generic Provisioning Control Format field values
	GenericProvisioningControlFormatFieldValues = map[Range]string{
		Range{0, 0}: "Transaction Start",
		Range{1, 1}: "Transaction Acknowledgment",
		Range{2, 2}: "Transaction Continuation",
		Range{3, 3}: "Provisioning Bearer Control",
	}

	// Table 5.17: Number of Elements field values
	NumberOfElementsFieldValues = map[Range]string{
		Range{0, 0}:   "Prohibited",
		Range{1, 255}: "The number of elements supported by the device",
	}

	// Table 5.21: Output OOB Size field values
	OutputOOBSizeFieldValues = map[Range]string{
		Range{0, 0}:   "The device does not support output OOB",
		Range{1, 8}:   "Maximum size in octets supported by the device",
		Range{9, 255}: "Reserved for Future Use",
	}

	// Table 5.23: Input OOB Size field values
	InputOOBSizeFieldValues = map[Range]string{
		Range{0, 0}:   "The device does not support Input OOB",
		Range{1, 8}:   "Maximum supported size in octets supported by the device",
		Range{9, 255}: "Reserved for Future Use",
	}

	// Table 5.26: Algorithm field values
	AlgorithmFieldValues = map[Range]string{
		Range{0, 0}:   "FIPS P-256 Elliptic Curve",
		Range{1, 255}: "Reserved for Future Use",
	}

	// Table 5.27: Public Key field values
	PublicKeyFieldValues = map[Range]string{
		Range{0, 0}:   "No OOB Public Key is used",
		Range{1, 1}:   "OOB Public Key is used",
		Range{2, 255}: "Prohibited",
	}

	// Table 5.28: Authentication Method field values
	AuthenticationMethodFieldValues = map[Range]string{
		Range{0, 0}:   "No OOB authentication is used",
		Range{1, 1}:   "Static OOB authentication is used",
		Range{2, 2}:   "Output OOB authentication is used",
		Range{3, 3}:   "Input OOB authentication is used",
		Range{4, 255}: "Prohibited",
	}

	// Table 5.29: Output OOB Action values for the Authentication Action field
	OutputOOBActionValuesForTheAuthenticationActionField = map[Range]string{
		Range{0, 0}:   "Blink",
		Range{1, 1}:   "Beep",
		Range{2, 2}:   "Vibrate",
		Range{3, 3}:   "Output Numeric",
		Range{4, 4}:   "Output Alphanumeric",
		Range{5, 255}: "Reserved for Future Use",
	}

	// Table 5.30: Output OOB Size values for the Authentication Size field
	OutputOOBSizeValuesForTheAuthenticationSizeField = map[Range]string{
		Range{0, 0}:   "Prohibited",
		Range{1, 8}:   "The Output OOB Size in characters to be used",
		Range{9, 255}: "Reserved for Future Use",
	}

	// Table 5.31: Input OOB Action values for the Authentication Action field
	InputOOBActionValuesForTheAuthenticationActionField = map[Range]string{
		Range{0, 0}:   "Push",
		Range{1, 1}:   "Twist",
		Range{2, 2}:   "Input Numeric",
		Range{3, 3}:   "Input Alphanumeric",
		Range{4, 255}: "Reserved for Future Use",
	}

	// Table 5.32: Input OOB Size values for the Authentication Size field
	InputOOBSizeValuesForTheAuthenticationSizeField = map[Range]string{
		Range{0, 0}:   "Prohibited",
		Range{1, 8}:   "The Input OOB size in characters to be used",
		Range{9, 255}: "Reserved for Future Use",
	}

	// Table 6.2: SAR field values
	SARFieldValues = map[Range]string{
		Range{0, 0}: "Data field contains a complete message",
		Range{1, 1}: "Data field contains the first segment of a message",
		Range{2, 2}: "Data field contains a continuation segment of a message",
		Range{3, 3}: "Data field contains the last segment of a message",
	}

	// Table 7.8: Identification Type values
	IdentificationTypeValues = map[Range]string{
		Range{0, 0}:   "Network ID type",
		Range{1, 1}:   "Node Identity type",
		Range{2, 255}: "Reserved for Future Use",
	}

	// Table 3.1: Generic OnOff states
	GenericOnOffStates = map[Range]string{
		Range{0, 0}:   "Off",
		Range{1, 1}:   "On",
		Range{2, 255}: "Prohibited",
	}

	// Table 3.2: Generic Level states
	GenericLevelStates = map[Range]string{
		Range{0, 65535}: "The Generic Level state of an element, represented as a 16-bit signed integer (the complement of 2)",
	}

	// Table 3.4: Default Transition Step Resolution values
	DefaultTransitionStepResolutionValues = map[Range]string{
		Range{0, 0}: "The Default Transition Step Resolution is 100 milliseconds",
		Range{1, 1}: "The Default Transition Step Resolution is 1 second",
		Range{2, 2}: "The Default Transition Step Resolution is 10 seconds",
		Range{3, 3}: "The Default Transition Step Resolution is 10 minutes",
	}

	// Table 3.5: Default Transition Number of Steps values
	DefaultTransitionNumberOfStepsValues = map[Range]string{
		Range{0, 0}:   "The Generic Default Transition Time is immediate.",
		Range{1, 62}:  "The number of steps.",
		Range{63, 63}: "The value is unknown. The state cannot be set to this value, but an element may report an unknown value if a transition is higher than 0x3E or not determined.",
	}

	// Table 3.6: Generic OnPowerUp states
	GenericOnPowerUpStates = map[Range]string{
		Range{0, 0}:   "Off. After being powered up, the element is in an off state.",
		Range{1, 1}:   "Default. After being powered up, the element is in an On state and uses default state values.",
		Range{2, 2}:   "Restore. If a transition was in progress when powered down, the element restores the target state when powered up. Otherwise the element restores the state it was in when powered down.",
		Range{3, 255}: "Prohibited",
	}

	// Table 3.7: Generic Power Actual states
	GenericPowerActualStates = map[Range]string{
		Range{0, 65535}: "Represents the power level relative to the maximum power level",
	}

	// Table 3.8: Generic Power Last states
	GenericPowerLastStates = map[Range]string{
		Range{0, 0}:     "Prohibited",
		Range{1, 65535}: "Represents the power level relative to the maximum power level",
	}

	// Table 3.9: Generic Power Default states
	GenericPowerDefaultStates = map[Range]string{
		Range{0, 0}:     "Use the Power Last value (see Section 3.1.5.1.1).",
		Range{1, 65535}: "Represents the power level relative to the maximum power level.",
	}

	// Table 3.10: Generic Power Min and Generic Power Max states
	GenericPowerMinAndGenericPowerMaxStates = map[Range]string{
		Range{0, 0}:     "Prohibited",
		Range{1, 65535}: "Represents the power level relative to the maximum power level.",
	}

	// Table 3.11: Generic Battery Level states
	GenericBatteryLevelStates = map[Range]string{
		Range{0, 100}:   "The percentage of the charge level. 100% represents fully charged. 0% represents fully discharged.",
		Range{101, 254}: "Prohibited",
		Range{255, 255}: "The percentage of the charge level is unknown.",
	}

	// Table 3.12: Generic Battery Time to Discharge states
	GenericBatteryTimeToDischargeStates = map[Range]string{
		Range{0, 16777214}:        "The remaining time (in minutes) of the discharging process",
		Range{16777215, 16777215}: "The remaining time of the discharging process is not known.",
	}

	// Table 3.13: Generic Battery Time to Charge states
	GenericBatteryTimeToChargeStates = map[Range]string{
		Range{0, 16777214}:        "The remaining time (in minutes) of the charging process",
		Range{16777215, 16777215}: "The remaining time of the charging process is not known.",
	}

	// Table 3.15: Generic Battery Flags Presence states
	GenericBatteryFlagsPresenceStates = map[Range]string{
		Range{0, 0}: "The battery is not present.",
		Range{1, 1}: "The battery is present and is removable.",
		Range{2, 2}: "The battery is present and is non-removable",
		Range{3, 3}: "The battery presence is unknown.",
	}

	// Table 3.16: Generic Battery Flags Indicator states
	GenericBatteryFlagsIndicatorStates = map[Range]string{
		Range{0, 0}: "The battery charge is Critically Low Level.",
		Range{1, 1}: "The battery charge is Low Level.",
		Range{2, 2}: "The battery charge is Good Level.",
		Range{3, 3}: "The battery charge is unknown.",
	}

	// Table 3.17: Generic Battery Flags Charging states
	GenericBatteryFlagsChargingStates = map[Range]string{
		Range{0, 0}: "The battery is not chargeable.",
		Range{1, 1}: "The battery is chargeable and is not charging.",
		Range{2, 2}: "The battery is chargeable and is charging.",
		Range{3, 3}: "The battery charging state is unknown.",
	}

	// Table 3.18: Generic Battery Flags Serviceability states
	GenericBatteryFlagsServiceabilityStates = map[Range]string{
		Range{0, 0}: "Reserved for Future Use",
		Range{1, 1}: "The battery does not require service.",
		Range{2, 2}: "The battery requires service.",
		Range{3, 3}: "The battery serviceability is unknown.",
	}

	// Table 4.3: Sensor Positive Tolerance states
	SensorPositiveToleranceStates = map[Range]string{
		Range{0, 0}:    "Unspecified",
		Range{1, 4095}: "The positive tolerance of the sensor. See Note below.",
	}

	// Table 4.4: Sensor Negative Tolerance states
	SensorNegativeToleranceStates = map[Range]string{
		Range{0, 0}:    "Unspecified",
		Range{1, 4095}: "The negative tolerance of the sensor. See Note below.",
	}

	// Table 4.5: Sensor sampling functions
	SensorSamplingFunctions = map[Range]string{
		Range{0, 0}:   "Unspecified",
		Range{1, 1}:   "Instantaneous",
		Range{2, 2}:   "Arithmetic Mean",
		Range{3, 3}:   "RMS",
		Range{4, 4}:   "Maximum",
		Range{5, 5}:   "Minimum",
		Range{6, 6}:   "Accumulated. (See note below.)",
		Range{7, 7}:   "Count. (See note below.)",
		Range{8, 255}: "Reserved for Future Use",
	}

	// Table 4.31: Sensor Data Format values
	SensorDataFormatValues = map[Range]string{
		Range{0, 0}: "Format A",
		Range{1, 1}: "Format B",
	}

	// Table 5.3: Scene Number values
	SceneNumberValues = map[Range]string{
		Range{0, 0}:     "Prohibited",
		Range{1, 65535}: "Scene Number value",
	}

	// Table 5.5: Year field values
	YearFieldValues = map[Range]string{
		Range{0, 99}:    "2 least significant digits of the year",
		Range{100, 100}: "Any year",
		Range{101, 127}: "Prohibited",
	}

	// Table 5.7: Day field values
	DayFieldValues = map[Range]string{
		Range{0, 0}:  "Any day",
		Range{1, 31}: "Day of the month",
	}

	// Table 5.8: Hour field values
	HourFieldValues = map[Range]string{
		Range{0, 23}:  "Hour of the day (00 to 23 hours)",
		Range{24, 24}: "Any hour of the day",
		Range{25, 25}: "Once a day (at a random hour)",
		Range{26, 31}: "Prohibited",
	}

	// Table 5.9: Minute field values
	MinuteFieldValues = map[Range]string{
		Range{0, 59}:  "Minute of the hour (00 to 59)",
		Range{60, 60}: "Any minute of the hour",
		Range{61, 61}: "Every 15 minutes (minute modulo 15 is 0) (0, 15, 30, 45)",
		Range{62, 62}: "Every 20 minutes (minute modulo 20 is 0) (0, 20, 40)",
		Range{63, 63}: "Once an hour (at a random minute)",
	}

	// Table 5.10: Second field values
	SecondFieldValues = map[Range]string{
		Range{0, 59}:  "Second of the minute (00 to 59)",
		Range{60, 60}: "Any second of the minute",
		Range{61, 61}: "Every 15 seconds (minute modulo 15 is 0) (0, 15, 30, 45)",
		Range{62, 62}: "Every 20 seconds (minute modulo 20 is 0) (0, 20, 40)",
		Range{63, 63}: "Once an minute (at a random second)",
	}

	// Table 5.12: Action field values
	ActionFieldValues = map[Range]string{
		Range{0, 0}:   "Turn Off",
		Range{1, 1}:   "Turn On",
		Range{2, 2}:   "Scene Recall",
		Range{15, 15}: "No action",
		Range{3, 14}:  "Reserved for Future Use",
	}

	// Table 5.13: Scene Number field values
	SceneNumberFieldValues = map[Range]string{
		Range{0, 0}:     "No scene",
		Range{1, 65535}: "Scene number",
	}

	// Table 5.31: Status code values
	StatusCodeValues = map[Range]string{
		Range{0, 0}:   "Success",
		Range{1, 1}:   "Scene Register Full",
		Range{2, 2}:   "Scene Not Found",
		Range{3, 255}: "Reserved for Future Use",
	}

	// Table 6.1: Light Lightness Linear states
	LightLightnessLinearStates = map[Range]string{
		Range{0, 0}:         "Light is not emitted by the element.",
		Range{1, 65534}:     "The lightness of a light emitted by the element.",
		Range{65535, 65535}: "The highest lightness of a light emitted by the element.",
	}

	// Table 6.2: Light Lightness Actual states
	LightLightnessActualStates = map[Range]string{
		Range{0, 0}:         "Light is not emitted by the element.",
		Range{1, 65534}:     "The perceived lightness of a light emitted by the element.",
		Range{65535, 65535}: "The highest perceived lightness of a light emitted by the element.",
	}

	// Table 6.3: Light Lightness Last states
	LightLightnessLastStates = map[Range]string{
		Range{0, 0}:         "Prohibited",
		Range{1, 65534}:     "The perceived lightness of a light emitted by the element",
		Range{65535, 65535}: "The highest perceived lightness of a light emitted by the element",
	}

	// Table 6.4: Light Lightness Default states
	LightLightnessDefaultStates = map[Range]string{
		Range{0, 0}:         "Use the Light Lightness Last value (see Section 6.1.2.3)",
		Range{1, 65534}:     "The perceived lightness of a light emitted by the element",
		Range{65535, 65535}: "The highest perceived lightness of a light emitted by the element",
	}

	// Table 6.5: Light Lightness Min and Light Lightness Max states
	LightLightnessMinAndLightLightnessMaxStates = map[Range]string{
		Range{0, 0}:     "Prohibited",
		Range{1, 65535}: "The lightness of an element",
	}

	// Table 6.6: Light CTL Temperature states
	LightCTLTemperatureStates = map[Range]string{
		Range{800, 20000}:   "The color temperature of white light in Kelvin",
		Range{0, 799}:       "Prohibited",
		Range{20001, 65535}: "Prohibited",
	}

	// Table 6.7: Light CTL Temperature Default states
	LightCTLTemperatureDefaultStates = map[Range]string{
		Range{800, 20000}:   "The color temperature of white light in Kelvin (0x0320 = 800 Kelvin, 0x4E20 = 20000 Kelvin)",
		Range{0, 799}:       "Prohibited",
		Range{20001, 65535}: "Prohibited",
	}

	// Table 6.8: Light CTL Temperature Range Min and Light CTL Temperature Range Max states
	LightCTLTemperatureRangeMinAndLightCTLTemperatureRangeMaxStates = map[Range]string{
		Range{800, 20000}:   "The color temperature of white light in Kelvin (0x0320 = 800 Kelvin, 0x4E20 = 20000 Kelvin)",
		Range{65535, 65535}: "The color temperature of white light is unknown",
		Range{0, 799}:       "Prohibited",
		Range{20001, 65534}: "Prohibited",
	}

	// Table 6.9: Light CTL Delta UV states
	LightCTLDeltaUVStates = map[Range]string{
		Range{32768, 32767}: "The 16-bit signed value representing the Delta UV of a tunable white light. A value of 0x0000 represents the Delta UV = 0 of a tunable white light.",
	}

	// Table 6.10: Light CTL Delta UV Default states
	LightCTLDeltaUVDefaultStates = map[Range]string{
		Range{32768, 32767}: "The 16-bit signed value representing the Delta UV of a tunable white light. A value of 0x0000 represents the Delta UV = 0 of a tunable white light.",
	}

	// Table 6.11: Light CTL Lightness states
	LightCTLLightnessStates = map[Range]string{
		Range{0, 0}:         "Light is not emitted by the element",
		Range{1, 65534}:     "The perceived lightness of a light emitted by the element",
		Range{65535, 65535}: "The highest perceived lightness of a light emitted by the element",
	}

	// Table 6.12: Light HSL Hue states
	LightHSLHueStates = map[Range]string{
		Range{0, 65535}: "The 16-bit value representing the hue",
	}

	// Table 6.13: Light HSL Hue Default states
	LightHSLHueDefaultStates = map[Range]string{
		Range{0, 65535}: "The 16-bit value representing the hue",
	}

	// Table 6.14: Light HSL Hue Min and Light HSL Hue Max states
	LightHSLHueMinAndLightHSLHueMaxStates = map[Range]string{
		Range{0, 65535}: "The hue of an element",
	}

	// Table 6.15: Light HSL Saturation states
	LightHSLSaturationStates = map[Range]string{
		Range{0, 0}:         "The lowest perceived saturation of a color light",
		Range{1, 65534}:     "The 16-bit value representing the saturation of a color light",
		Range{65535, 65535}: "The highest perceived saturation of a color light",
	}

	// Table 6.16: Light HSL Saturation Default states
	LightHSLSaturationDefaultStates = map[Range]string{
		Range{0, 65535}: "The 16-bit value representing the saturation",
	}

	// Table 6.17: Light HSL Saturation Min and Light HSL Saturation Max states
	LightHSLSaturationMinAndLightHSLSaturationMaxStates = map[Range]string{
		Range{0, 65535}: "The saturation of an element",
	}

	// Table 6.18: Light HSL Lightness states
	LightHSLLightnessStates = map[Range]string{
		Range{0, 0}:         "Light is not emitted by the element",
		Range{1, 65534}:     "The perceived lightness of a light emitted by the element",
		Range{65535, 65535}: "The highest perceived lightness of a light emitted by the element",
	}

	// Table 6.19: Light xyL x states
	LightXyLXStates = map[Range]string{
		Range{0, 0}:         "The value of 0 representing the x coordinate of a CIE1931 color light",
		Range{1, 65534}:     "The 16-bit value representing the x coordinate of a CIE1931 color light",
		Range{65535, 65535}: "The value of 1 representing the x coordinate of a CIE1931 color light",
	}

	// Table 6.20: Light xyL x Default states
	LightXyLXDefaultStates = map[Range]string{
		Range{0, 65535}: "The 16-bit value representing the x",
	}

	// Table 6.21: Light xyL x Min and Light xyL x Max states
	LightXyLXMinAndLightXyLXMaxStates = map[Range]string{
		Range{0, 65535}: "The value of a Light xyL x state of an element",
	}

	// Table 6.22: Light xyL y states
	LightXyLYStates = map[Range]string{
		Range{0, 0}:         "The value of 0 representing the y coordinate of a CIE1931 color light",
		Range{1, 65534}:     "The 16-bit value representing the y coordinate of a CIE1931 color light",
		Range{65535, 65535}: "The value of 1 representing the y coordinate of a CIE1931 color light",
	}

	// Table 6.23: Light xyL y Default states
	LightXyLYDefaultStates = map[Range]string{
		Range{0, 65535}: "The 16-bit value representing the y",
	}

	// Table 6.24: Light xyL y Min and Light xyL y Max states
	LightXyLYMinAndLightXyLYMaxStates = map[Range]string{
		Range{0, 65535}: "The value of a Light xyL y state of an element",
	}

	// Table 6.25: Light xyL Lightness states
	LightXyLLightnessStates = map[Range]string{
		Range{0, 0}:         "Light is not emitted by the element",
		Range{1, 65534}:     "The perceived lightness of a light emitted by the element",
		Range{65535, 65535}: "The highest perceived lightness of a light emitted by the element.",
	}

	// Table 6.26: Light LC Mode states
	LightLCModeStates = map[Range]string{
		Range{0, 0}: "The controller is turned off. The binding with the Light Lightness state is disabled.",
		Range{1, 1}: "The controller is turned on. The binding with the Light Lightness state is enabled.",
	}

	// Table 6.27: Light LC Occupancy Mode states
	LightLCOccupancyModeStates = map[Range]string{
		Range{0, 0}: "The controller does not transition from a standby state when occupancy is reported.",
		Range{1, 1}: "The controller may transition from a standby state when occupancy is reported.",
	}

	// Table 6.28: Light LC Light OnOff states
	LightLCLightOnOffStates = map[Range]string{
		Range{0, 0}: "Off or Standby",
		Range{1, 1}: "Occupancy or Run or Prolong",
	}

	// Table 6.29: Light LC Occupancy states
	LightLCOccupancyStates = map[Range]string{
		Range{0, 0}: "There is no occupancy reported by occupancy sensors.",
		Range{1, 1}: "There has been occupancy reported by occupancy sensors.",
	}

	// Table 6.30: Light LC Ambient LuxLevel states
	LightLCAmbientLuxLevelStates = map[Range]string{
		Range{0, 16777215}: "Illuminance from 0.00 to 167772.16 lux",
	}

	// Table 6.31: Light LC Linear Output states
	LightLCLinearOutputStates = map[Range]string{
		Range{0, 0}:         "Light is not emitted by the element.",
		Range{1, 65534}:     "The lightness of a light emitted by the element.",
		Range{65535, 65535}: "The highest lightness of a light emitted by the element.",
	}

	// Table 6.32: Light LC Ambient LuxLevel On states
	LightLCAmbientLuxLevelOnStates = map[Range]string{
		Range{0, 65535}: "Illuminance from 0 to 65535 lux",
	}

	// Table 6.33: Light LC Ambient LuxLevel Prolong states
	LightLCAmbientLuxLevelProlongStates = map[Range]string{
		Range{0, 65535}: "Illuminance from 0 to 65535 lux",
	}

	// Table 6.34: Light LC Ambient LuxLevel Standby states
	LightLCAmbientLuxLevelStandbyStates = map[Range]string{
		Range{0, 65535}: "Illuminance from 0 to 65535 lux",
	}
)
