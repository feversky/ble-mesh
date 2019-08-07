package def

type (

	// Table 3.3: Generic Default Transition Time state format
	GenericDefaultTransitionTimeStateFormat struct {
		DefaultTransitionNumberOfSteps  uint `bits:"6"` // The number of Steps
		DefaultTransitionStepResolution uint `bits:"2"` // The resolution of the Default Transition Number of Steps field
	}

	// Table 3.19: Generic Location state
	GenericLocationState struct {
		GlobalLatitude  uint `bits:"32"` // Global Coordinates (Latitude)
		GlobalLongitude uint `bits:"32"` // Global Coordinates (Longitude)
		GlobalAltitude  uint `bits:"16"` // Global Altitude
		LocalNorth      uint `bits:"16"` // Local Coordinates (North)
		LocalEast       uint `bits:"16"` // Local Coordinates (East)
		LocalAltitude   uint `bits:"16"` // Local Altitude
		FloorNumber     uint `bits:"8"`  // Floor Number
		Uncertainty     uint `bits:"16"` // Uncertainty
	}

	// Table 3.24: Generic User Property states
	GenericUserPropertyStates struct {
		UserPropertyID    uint `bits:"16"`       // Defined in Section 3.1.8.1.1.
		UserAccess        uint `bits:"8"`        // Defined in Section 3.1.8.1.2.
		UserPropertyValue uint `bits:"variable"` // Scalar or String value, defined in Section 3.1.8.1.3.
	}

	// Table 3.27: Generic Admin Property states
	GenericAdminPropertyStates struct {
		AdminPropertyID    uint `bits:"16"`       // Defined in Section 3.1.8.2.1.
		AdminUserAccess    uint `bits:"8"`        // Defined in Section 3.1.8.2.2.
		AdminPropertyValue uint `bits:"variable"` // Scalar or String value, defined in Section 3.1.8.2.3.
	}

	// Table 3.30: Generic Manufacturer Property states
	GenericManufacturerPropertyStates struct {
		ManufacturerPropertyID    uint `bits:"16"`       // Defined in Section 3.1.8.3.1.
		ManufacturerUserAccess    uint `bits:"8"`        // Defined in Section 3.1.8.3.2.
		ManufacturerPropertyValue uint `bits:"variable"` // Scalar or String value, defined in Section 3.1.8.3.3.
	}

	// Table 3.33: Generic Client Property state
	GenericClientPropertyState struct {
		ClientPropertyID uint `bits:"16"` // Defined in Section 3.1.9.1.
	}

	// Table 3.35: Generic OnOff Set message parameters
	GenericOnOffSetMessageParameters struct {
		OnOff          uint                                    `bits:"8"` // The target value of the Generic OnOff state
		TID            uint                                    `bits:"8"` // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"` // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"` // Message execution delay in 5 millisecond steps (C.1)
	}

	// Table 3.36: Generic OnOff Set Unacknowledged message parameters
	GenericOnOffSetUnacknowledgedMessageParameters struct {
		OnOff          uint                                    `bits:"8"` // The target value of the Generic OnOff state
		TID            uint                                    `bits:"8"` // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"` // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"` // Message execution delay in 5 milliseconds steps (so)
	}

	// Table 3.37: Generic OnOff Status message parameters
	GenericOnOffStatusMessageParameters struct {
		PresentOnOff  uint `bits:"8"`  // The present value of the Generic OnOff state.
		TargetOnOff   uint `bits:"t8"` // The target value of the Generic OnOff state (optional).
		RemainingTime uint `bits:"8"`  // Format as defined in Section 3.1.3. (C.1)
	}

	// Table 3.38: Generic Level Set message parameters
	GenericLevelSetMessageParameters struct {
		Level          uint                                    `bits:"16"` // The target value of the Generic Level state
		TID            uint                                    `bits:"8"`  // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"`  // Message execution delay in 5 milliseconds steps (C.1)
	}

	// Table 3.39: Generic Level Set Unacknowledged message parameters
	GenericLevelSetUnacknowledgedMessageParameters struct {
		Level          uint                                    `bits:"16"` // The target value of the Generic Level state
		TID            uint                                    `bits:"8"`  // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"`  // Message execution delay in 5 millisecond steps (C.1)
	}

	// Table 3.40: Generic Delta Set message parameters
	GenericDeltaSetMessageParameters struct {
		DeltaLevel     uint                                    `bits:"32"` // The Delta change of the Generic Level state
		TID            uint                                    `bits:"8"`  // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"`  // Message execution delay in 5 milliseconds steps (C.1)
	}

	// Table 3.41: Generic Delta Set Unacknowledged message parameters
	GenericDeltaSetUnacknowledgedMessageParameters struct {
		DeltaLevel     uint                                    `bits:"32"` // The Delta change of the Generic Level state
		TID            uint                                    `bits:"8"`  // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"`  // Message execution delay in 5 milliseconds steps (C.1)
	}

	// Table 3.42: Generic Move Set message parameters
	GenericMoveSetMessageParameters struct {
		DeltaLevel     uint                                    `bits:"16"` // The Delta Level step to calculate Move speed for the Generic Level state.
		TID            uint                                    `bits:"8"`  // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"`  // Format as defined in Section 3.1.3 (optional).
		Delay          uint                                    `bits:"8"`  // Message execution delay in 5 milliseconds steps (C.1).
	}

	// Table 3.43: Generic Move Set Unacknowledged message parameters
	GenericMoveSetUnacknowledgedMessageParameters struct {
		DeltaLevel     uint                                    `bits:"16"` // The Delta Level step to calculate Move speed for the Generic Level state
		TID            uint                                    `bits:"8"`  // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"`  // Format as defined in Section 3.1.3. (Optional).
		Delay          uint                                    `bits:"8"`  // Message execution delay in 5 milliseconds steps (C.1).
	}

	// Table 3.44: Generic Level Status message parameters
	GenericLevelStatusMessageParameters struct {
		PresentLevel  uint `bits:"16"`  // The present value of the Generic Level state.
		TargetLevel   uint `bits:"t16"` // The target value of the Generic Level state (Optional).
		RemainingTime uint `bits:"8"`   // Format as defined in Section 3.1.3 (C.1).
	}

	// Table 3.45: Generic Default Transition Time Set message parameters
	GenericDefaultTransitionTimeSetMessageParameters struct {
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"` // The value of the Generic Default Transition Time state.
	}

	// Table 3.46: Generic Default Transition Time Set Unacknowledged message parameters
	GenericDefaultTransitionTimeSetUnacknowledgedMessageParameters struct {
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"` // The value of the Generic Default Transition Time state.
	}

	// Table 3.47: Generic Default Transition Time Status message parameters
	GenericDefaultTransitionTimeStatusMessageParameters struct {
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"` // The value of the Generic Default Transition Time state.
	}

	// Table 3.48: Generic OnPowerUp Set message parameters
	GenericOnPowerUpSetMessageParameters struct {
		OnPowerUp uint `bits:"8"` // The value of the Generic OnPowerUp state.
	}

	// Table 3.49: Generic OnPowerUp Set Unacknowledged message parameters
	GenericOnPowerUpSetUnacknowledgedMessageParameters struct {
		OnPowerUp uint `bits:"8"` // The value of the Generic OnPowerUp state.
	}

	// Table 3.50: Generic OnPowerUp Status message parameters
	GenericOnPowerUpStatusMessageParameters struct {
		OnPowerUp uint `bits:"8"` // The value of the Generic OnPowerUp state.
	}

	// Table 3.51: Generic Power Level Set message parameters
	GenericPowerLevelSetMessageParameters struct {
		Power          uint                                    `bits:"16"` // The target value of the Generic Power Actual state
		TID            uint                                    `bits:"8"`  // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"`  // Message execution delay in 5 milliseconds steps (C.1)
	}

	// Table 3.52: Generic Power Level Set Unacknowledged message parameters
	GenericPowerLevelSetUnacknowledgedMessageParameters struct {
		Power          uint                                    `bits:"16"` // The target value of the Generic Power Actual state
		TID            uint                                    `bits:"8"`  // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"`  // Message execution delay in 5 milliseconds steps (C.1)
	}

	// Table 3.53: Generic Power Level Status message parameters
	GenericPowerLevelStatusMessageParameters struct {
		PresentPower  uint `bits:"16"`  // The present value of the Generic Power Actual state.
		TargetPower   uint `bits:"t16"` // The target value of the Generic Power Actual state (optional).
		RemainingTime uint `bits:"8"`   // Format as defined in Section 3.1.3 (C.1).
	}

	// Table 3.54: Generic Power Last Status message parameters
	GenericPowerLastStatusMessageParameters struct {
		Power uint `bits:"16"` // The value of the Generic Power Last state.
	}

	// Table 3.55: Generic Power Default Set message parameters
	GenericPowerDefaultSetMessageParameters struct {
		Power uint `bits:"16"` // The value of the Generic Power Default state.
	}

	// Table 3.56: Generic Power Default Set Unacknowledged message parameters
	GenericPowerDefaultSetUnacknowledgedMessageParameters struct {
		Power uint `bits:"16"` // The value of the Generic Power Default state.
	}

	// Table 3.57: Generic Power Default Status message parameters
	GenericPowerDefaultStatusMessageParameters struct {
		Power uint `bits:"16"` // The value of the Generic Power Default state.
	}

	// Table 3.58: Generic Power Range Set message parameters
	GenericPowerRangeSetMessageParameters struct {
		RangeMin uint `bits:"16"` // The value of the Generic Power Min field of the Generic Power Range state.
		RangeMax uint `bits:"16"` // The value of the Generic Power Range Max field of the Generic Power Range state.
	}

	// Table 3.59: Generic Power Range Set Unacknowledged message parameters
	GenericPowerRangeSetUnacknowledgedMessageParameters struct {
		RangeMin uint `bits:"16"` // The value of the Generic Power Min field of the Generic Power Range state.
		RangeMax uint `bits:"16"` // The value of the Generic Power Range Max field of the Generic Power Range state.
	}

	// Table 3.60: Generic Power Range Status message parameters
	GenericPowerRangeStatusMessageParameters struct {
		StatusCode uint `bits:"8"`  // Status Code for the requesting message.
		RangeMin   uint `bits:"16"` // The value of the Generic Power Range Min field of the Generic Power Range state.
		RangeMax   uint `bits:"16"` // The value of the Generic Power Range Max field of the Generic Power Range state.
	}

	// Table 3.61: Generic Battery Status message parameters
	GenericBatteryStatusMessageParameters struct {
		BatteryLevel    uint `bits:"8"`  // The value of the Generic Battery Level state.
		TimeToDischarge uint `bits:"24"` // The value of the Generic Battery Time to Discharge state.
		TimeToCharge    uint `bits:"24"` // The value of the Generic Battery Time to Charge state.
		Flags           uint `bits:"8"`  // The value of the Generic Battery Flags state.
	}

	// Table 3.62: Generic Location Global Set message parameters
	GenericLocationGlobalSetMessageParameters struct {
		GlobalLatitude  uint `bits:"32"` // Global Coordinates (Latitude)
		GlobalLongitude uint `bits:"32"` // Global Coordinates (Longitude)
		GlobalAltitude  uint `bits:"16"` // Global Altitude
	}

	// Table 3.63: Generic Location Global Set Unacknowledged message parameters
	GenericLocationGlobalSetUnacknowledgedMessageParameters struct {
		GlobalLatitude  uint `bits:"32"` // Global Coordinates (Latitude)
		GlobalLongitude uint `bits:"32"` // Global Coordinates (Longitude)
		GlobalAltitude  uint `bits:"16"` // Global Altitude
	}

	// Table 3.64: Generic Location Global Status message parameters
	GenericLocationGlobalStatusMessageParameters struct {
		GlobalLatitude  uint `bits:"32"` // Global Coordinates (Latitude)
		GlobalLongitude uint `bits:"32"` // Global Coordinates (Longitude)
		GlobalAltitude  uint `bits:"16"` // Global Altitude
	}

	// Table 3.65: Generic Location Local Set message parameters
	GenericLocationLocalSetMessageParameters struct {
		LocalNorth    uint `bits:"16"` // Local Coordinates (North)
		LocalEast     uint `bits:"16"` // Local Coordinates (East)
		LocalAltitude uint `bits:"16"` // Local Altitude
		FloorNumber   uint `bits:"8"`  // Floor Number
		Uncertainty   uint `bits:"16"` // Uncertainty
	}

	// Table 3.66: Generic Location Local Set Unacknowledged message parameters
	GenericLocationLocalSetUnacknowledgedMessageParameters struct {
		LocalNorth    uint `bits:"16"` // Local Coordinates (North)
		LocalEast     uint `bits:"16"` // Local Coordinates (East)
		LocalAltitude uint `bits:"16"` // Local Altitude
		FloorNumber   uint `bits:"8"`  // Floor Number
		Uncertainty   uint `bits:"16"` // Uncertainty
	}

	// Table 3.67: Generic Location Local Status message parameters
	GenericLocationLocalStatusMessageParameters struct {
		LocalNorth    uint `bits:"16"` // Local Coordinates (North)
		LocalEast     uint `bits:"16"` // Local Coordinates (East)
		LocalAltitude uint `bits:"16"` // Local Altitude
		FloorNumber   uint `bits:"8"`  // Floor Number
		Uncertainty   uint `bits:"16"` // Uncertainty
	}

	// Table 3.68: Generic User Properties Status message parameters
	GenericUserPropertiesStatusMessageParameters struct {
		UserPropertyIDs uint `bits:"2*N"` // A sequence of N User Property IDs present within an element, where N is the number of device property IDs included in the message.
	}

	// Table 3.69: Generic User Property Get message parameters
	GenericUserPropertyGetMessageParameters struct {
		UserPropertyID uint `bits:"16"` // Property ID identifying a Generic User Property
	}

	// Table 3.70: Generic User Property Set message parameters
	GenericUserPropertySetMessageParameters struct {
		UserPropertyID    uint `bits:"16"`       // Property ID identifying a Generic User Property
		UserPropertyValue uint `bits:"variable"` // Raw value for the User Property
	}

	// Table 3.71: Generic User Property Set Unacknowledged message parameters
	GenericUserPropertySetUnacknowledgedMessageParameters struct {
		UserPropertyID    uint `bits:"16"`       // Property ID identifying a Generic User Property
		UserPropertyValue uint `bits:"variable"` // Raw value for the User Property
	}

	// Table 3.72: Generic User Property Status message parameters
	GenericUserPropertyStatusMessageParameters struct {
		UserPropertyID    uint `bits:"16"`       // Property ID identifying a Generic User Property.
		UserAccess        uint `bits:"t8"`       // Enumeration indicating user access (Optional)
		UserPropertyValue uint `bits:"variable"` // Raw value for the User Property (C.1)
	}

	// Table 3.73: Generic Admin Properties Status message parameters
	GenericAdminPropertiesStatusMessageParameters struct {
		AdminPropertyIDs uint `bits:"2*N"` // A sequence of N Admin Property IDs present within an element, where N is the number of device property IDs included in the message.
	}

	// Table 3.74: Generic Admin Property Get message parameters
	GenericAdminPropertyGetMessageParameters struct {
		AdminPropertyID uint `bits:"16"` // Property ID identifying a Generic Admin Property.
	}

	// Table 3.75: Generic Admin Property Set message parameters
	GenericAdminPropertySetMessageParameters struct {
		AdminPropertyID    uint `bits:"16"`       // Property ID identifying a Generic Admin Property.
		AdminUserAccess    uint `bits:"8"`        // Enumeration indicating user access.
		AdminPropertyValue uint `bits:"variable"` // Raw value for the Admin Property
	}

	// Table 3.76: Generic Admin Property Set Unacknowledged message parameters
	GenericAdminPropertySetUnacknowledgedMessageParameters struct {
		AdminPropertyID    uint `bits:"16"`       // Property ID identifying a Generic Admin Property.
		AdminUserAccess    uint `bits:"8"`        // Enumeration indicating user access.
		AdminPropertyValue uint `bits:"variable"` // Raw value for the Admin Property.
	}

	// Table 3.77: Generic Admin Property Status message parameters
	GenericAdminPropertyStatusMessageParameters struct {
		AdminPropertyID    uint `bits:"16"`       // Property ID identifying a Generic Admin Property
		AdminUserAccess    uint `bits:"t8"`       // Enumeration indicating user access (Optional)
		AdminPropertyValue uint `bits:"variable"` // Raw value for the Admin Property (C.1)
	}

	// Table 3.78: Generic Manufacturer Properties Status message parameters
	GenericManufacturerPropertiesStatusMessageParameters struct {
		ManufacturerPropertyIDs uint `bits:"2*N"` // A sequence of N Manufacturer Property IDs present within an element, where N is the number of device property IDs included in the message.
	}

	// Table 3.79: Generic Manufacturer Property Get message parameters
	GenericManufacturerPropertyGetMessageParameters struct {
		ManufacturerPropertyID uint `bits:"16"` // Property ID identifying a Generic Manufacturer Property
	}

	// Table 3.80: Generic Manufacturer Property Set message parameters
	GenericManufacturerPropertySetMessageParameters struct {
		ManufacturerPropertyID uint `bits:"16"` // Property ID identifying a Generic Manufacturer Property
		ManufacturerUserAccess uint `bits:"8"`  // Enumeration indicating user access
	}

	// Table 3.81: Generic Manufacturer Property Set Unacknowledged message parameters
	GenericManufacturerPropertySetUnacknowledgedMessageParameters struct {
		ManufacturerPropertyID uint `bits:"16"` // Property ID identifying a Generic Manufacturer Property
		ManufacturerUserAccess uint `bits:"8"`  // Enumeration indicating user access
	}

	// Table 3.82: Generic Manufacturer Property Status message parameters
	GenericManufacturerPropertyStatusMessageParameters struct {
		ManufacturerPropertyID    uint `bits:"16"`       // Property ID identifying a Generic Manufacturer Property
		ManufacturerUserAccess    uint `bits:"t8"`       // Enumeration indicating user access (Optional)
		ManufacturerPropertyValue uint `bits:"variable"` // Raw value for the Manufacturer Property (C.1)
	}

	// Table 3.83: Generic Client Properties Get message parameters
	GenericClientPropertiesGetMessageParameters struct {
		ClientPropertyID uint `bits:"16"` // A starting Client Property ID present within an element
	}

	// Table 3.84: Generic Client Properties Status message parameters
	GenericClientPropertiesStatusMessageParameters struct {
		ClientPropertyIDs uint `bits:"2*N"` // A sequence of N Client Property IDs present within an element, where N is the number of device property IDs included in the message.
	}

	// Table 4.1: Sensor Descriptor states
	SensorDescriptorStates struct {
		SensorPropertyID        uint `bits:"16"` // Defined in Section 4.1.1.1.
		SensorPositiveTolerance uint `bits:"12"` // Defined in Section 4.1.1.2.
		SensorNegativeTolerance uint `bits:"12"` // Defined in Section 4.1.1.3.
		SensorSamplingFunction  uint `bits:"8"`  // Defined in Section 4.1.1.4.
		SensorMeasurementPeriod uint `bits:"8"`  // Defined in Section 4.1.1.5.
		SensorUpdateInterval    uint `bits:"8"`  // Defined in Section 4.1.1.6.
	}

	// Table 4.8: Sensor Setting state
	SensorSettingState struct {
		SensorPropertyID        uint `bits:"16"`       // Property ID of a sensor
		SensorSettingPropertyID uint `bits:"16"`       // Property ID of a setting within a sensor
		SensorSettingAccess     uint `bits:"8"`        // Read/Write access rights for the setting
		SensorSettingRaw        uint `bits:"variable"` // Raw value of a setting within a sensor
	}

	// Table 4.12: Sensor Cadence states
	SensorCadenceStates struct {
		SensorProperty           uint `bits:"16"`       // Defined in Section 4.1.1.1.
		FastCadencePeriodDivisor uint `bits:"7"`        // Divisor for the Publish Period (see Mesh Profile specification [2]).
		StatusTriggerType        uint `bits:"1"`        // Defines the unit and format of the Status Trigger Delta fields.
		StatusTriggerDeltaDown   uint `bits:"variable"` // Delta down value that triggers a status message.
		StatusTriggerDeltaUp     uint `bits:"variable"` // Delta up value that triggers a status message.
		StatusMinInterval        uint `bits:"8"`        // Minimum interval between two consecutive Status messages.
		FastCadenceLow           uint `bits:"variable"` // Low value for the fast cadence range.
		FastCadenceHigh          uint `bits:"variable"` // High value for the fast cadence range.
	}

	// Table 4.13: Sensor Data state
	SensorDataState struct {
		Data []uint `bits:"8"` // ID of the nth device property of the sensor and Raw Value field with a size and representation defined by the nth device property
	}

	// Table 4.14: Sensor Series Column states
	SensorSeriesColumnStates struct {
		SensorPropertyID  uint `bits:"16"`       // Property describing the data series of the sensor
		SensorRawValueX   uint `bits:"variable"` // Raw value representing the left corner of a column on the X axis
		SensorColumnWidth uint `bits:"variable"` // Raw value representing the width of the column
		SensorRawValueY   uint `bits:"variable"` // Raw value representing the height of the column on the Y axis
	}

	// Table 4.16: Sensor Descriptor Get message parameters
	SensorDescriptorGetMessageParameters struct {
		PropertyID uint `bits:"16"` // Property ID for the sensor (Optional)
	}

	// Table 4.17: Sensor Descriptor Status message parameters
	SensorDescriptorStatusMessageParameters struct {
		Descriptor uint `bits:"8*N or 2"` // Sequence of 8 octet Sensor Descriptors (Optional)
	}

	// Table 4.18: Sensor Cadence Get message parameters
	SensorCadenceGetMessageParameters struct {
		PropertyID uint `bits:"16"` // Property ID for the sensor.
	}

	// Table 4.19: Sensor Cadence Set message parameters
	SensorCadenceSetMessageParameters struct {
		PropertyID               uint `bits:"16"`       // Property ID for the sensor.
		FastCadencePeriodDivisor uint `bits:"7"`        // Divisor for the Publish Period (see Mesh Profile specification [2]).
		StatusTriggerType        uint `bits:"1"`        // Defines the unit and format of the Status Trigger Delta fields.
		StatusTriggerDeltaDown   uint `bits:"variable"` // Delta down value that triggers a status message.
		StatusTriggerDeltaUp     uint `bits:"variable"` // Delta up value that triggers a status message.
		StatusMinInterval        uint `bits:"8"`        // Minimum interval between two consecutive Status messages.
		FastCadenceLow           uint `bits:"variable"` // Low value for the fast cadence range.
		FastCadenceHigh          uint `bits:"variable"` // High value for the fast cadence range.
	}

	// Table 4.20: Sensor Cadence Set Unacknowledged message parameters
	SensorCadenceSetUnacknowledgedMessageParameters struct {
		PropertyID               uint `bits:"16"`       // Property for the sensor.
		FastCadencePeriodDivisor uint `bits:"7"`        // Divisor for the Publish Period (see Mesh Profile specification [2]).
		StatusTriggerType        uint `bits:"1"`        // Defines the unit and format of the Status Trigger Delta fields.
		StatusTriggerDeltaDown   uint `bits:"variable"` // Delta down value that triggers a status message.
		StatusTriggerDeltaUp     uint `bits:"variable"` // Delta up value that triggers a status message.
		StatusMinInterval        uint `bits:"8"`        // Minimum interval between two consecutive Status messages.
		FastCadenceLow           uint `bits:"variable"` // Low value for the fast cadence range.
		FastCadenceHigh          uint `bits:"variable"` // High value for the fast cadence range.
	}

	// C.1: If the Fast Cadence Period Divisor field is present, the Status Trigger Type, Status Trigger Delta Down, Status Trigger Delta Up, Status Min Interval, Fast Cadence Low, and Fast Cadence High fields shall also be present; otherwise these fields shall not be present. Table 4.21: Sensor Cadence Status message parameters
	SensorCadenceStatus struct {
		PropertyID               uint `bits:"16"`       // Property for the sensor.
		FastCadencePeriodDivisor uint `bits:"7"`        // Divisor for the Publish Period (see Mesh Profile specification [2]). (Optional)
		StatusTriggerType        uint `bits:"1"`        // Defines the unit and format of the Status Trigger Delta fields. (C.1)
		StatusTriggerDeltaDown   uint `bits:"variable"` // Delta down value that triggers a status message. (C.1)
		StatusTriggerDeltaUp     uint `bits:"variable"` // Delta up value that triggers a status message. (C.1)
		StatusMinInterval        uint `bits:"8"`        // Minimum interval between two consecutive status messages. (C.1)
		FastCadenceLow           uint `bits:"variable"` // Low value for the fast cadence range. (C.1)
		FastCadenceHigh          uint `bits:"variable"` // High value for the fast cadence range. (C.1)
	}

	// Table 4.22: Sensor Settings Get message parameters
	SensorSettingsGetMessageParameters struct {
		SensorPropertyID uint `bits:"16"` // Property ID identifying a sensor.
	}

	// Table 4.23: Sensor Setting Status message parameters
	SensorSettingsStatusMessageParameters struct {
		SensorPropertyID         uint `bits:"16"`  // Property ID identifying a sensor.
		SensorSettingPropertyIDs uint `bits:"2*N"` // A sequence of N Sensor Setting Property IDs identifying settings within a sensor, where N is the number of property IDs included in the message. (Optional)
	}

	// Table 4.24: Sensor Setting Get message parameters
	SensorSettingGetMessageParameters struct {
		SensorPropertyID        uint `bits:"16"` // Property ID identifying a sensor.
		SensorSettingPropertyID uint `bits:"16"` // Setting Property ID identifying a setting within a sensor.
	}

	// Table 4.25: Sensor Setting Set message parameters
	SensorSettingSetMessageParameters struct {
		SensorPropertyID        uint `bits:"16"`       // Property ID identifying a sensor.
		SensorSettingPropertyID uint `bits:"16"`       // Setting ID identifying a setting within a sensor.
		SensorSettingRaw        uint `bits:"variable"` // Raw value for the setting.
	}

	// Table 4.26: Sensor Setting Set Unacknowledged message parameters
	SensorSettingSetUnacknowledgedMessageParameters struct {
		SensorPropertyID        uint `bits:"16"`       // Property ID identifying a sensor.
		SensorSettingPropertyID uint `bits:"16"`       // Setting ID identifying a setting within a sensor.
		SensorSettingRaw        uint `bits:"variable"` // Raw value for the setting.
	}

	// Table 4.27: Sensor Setting Status message parameters
	SensorSettingStatusMessageParameters struct {
		SensorPropertyID        uint `bits:"16"`       // Property ID identifying a sensor.
		SensorSettingPropertyID uint `bits:"16"`       // Setting ID identifying a setting within a sensor.
		SensorSettingAccess     uint `bits:"8"`        // Read / Write access rights for the setting. (Optional)
		SensorSettingRaw        uint `bits:"variable"` // Raw value for the setting. (C.1)
	}

	// Table 4.28: Sensor Get message parameters
	SensorGetMessageParameters struct {
		PropertyID uint `bits:"16"` // Property for the sensor. (Optional)
	}

	// Table 4.29: Sensor Status message parameters
	SensorStatusMessageParameters struct {
		MarshalledSensorData uint `bits:"variable"` // The Sensor Data state. (Optional)
	}

	// Table 4.30: Marshalled Sensor Data field
	MarshalledSensorDataField struct {
		MPID1     uint `bits:"2 or 3"`   // TLV of the 1st device property of the sensor.
		RawValue1 uint `bits:"variable"` // Raw Value field with a size and representation defined by the 1st device property.
		MPID2     uint `bits:"2 or 3"`   // TLV of the 2nd device property of a sensor.
		RawValue2 uint `bits:"variable"` // Raw Value field with a size and representation defined by the 2nd device property.
		MPIDN     uint `bits:"2 or 3"`   // TLV of the nth device property of the sensor.
		RawValueN uint `bits:"variable"` // Raw Value field with a size and representation defined by the nth device property.
	}

	// Table 4.32: Format A of the Marshalled Property ID (MPID) field The Format field is 0b0 and indicates that Format A is used.
	FormatAOfTheMarshalledPropertyID struct {
		Format     uint `bits:"1"`  // Format A tag, 0b0
		Length     uint `bits:"4"`  // Length of the Property Value
		PropertyID uint `bits:"11"` // Property identifying a sensor (Optional)
	}

	// Table 4.33: Format B of the Marshalled Property ID (MPID) field The Format field is 0b1 and indicates Format B is used.
	FormatBOfTheMarshalledPropertyID struct {
		Format     uint `bits:"1"`  // Format B tag, 0b1
		Length     uint `bits:"7"`  // Length of the Property Value
		PropertyID uint `bits:"16"` // Property identifying a sensor (Optional)
	}

	// Table 4.34: Sensor Column Get message parameters
	SensorColumnGetMessageParameters struct {
		PropertyID uint `bits:"16"`       // Property identifying a sensor
		RawValueX  uint `bits:"variable"` // Raw value identifying a column
	}

	// Table 4.35: Sensor Column Status message parameters
	SensorColumnStatusMessageParameters struct {
		PropertyID  uint `bits:"16"`       // Property identifying a sensor and the Y axis.
		RawValueX   uint `bits:"variable"` // Raw value representing the left corner of the column on the X axis.
		ColumnWidth uint `bits:"variable"` // Raw value representing the width of the column. (Optional)
		RawValueY   uint `bits:"variable"` // Raw value representing the height of the column on the Y axis. (C.1)
	}

	// Table 4.36: Sensor Series Get message parameters
	SensorSeriesGetMessageParameters struct {
		PropertyID uint `bits:"16"`       // Property identifying a sensor.
		RawValueX1 uint `bits:"variable"` // Raw value identifying a starting column. (Optional)
		RawValueX2 uint `bits:"variable"` // Raw value identifying an ending column. (C.1)
	}

	// Table 4.37: Sensor Series Status message parameters
	SensorSeriesStatusMessageParameters struct {
		PropertyID  uint   `bits:"16"`       // Property identifying a sensor and the Y axis.
		RawValueX   []uint `bits:"variable"` // Raw value representing the left corner of the nth column on the X axis.
		ColumnWidth []uint `bits:"variable"` // Raw value representing the width of the nth column.
		RawValueY   []uint `bits:"variable"` // Raw value representing the height of the nth column on the Y axis.
	}

	// Table 5.15: Time Set message parameters
	TimeSetMessageParameters struct {
		TAISeconds     uint `bits:"40"` // The current TAI time in seconds
		Subsecond      uint `bits:"8"`  // The sub second time in units of 1/256th second
		Uncertainty    uint `bits:"8"`  // The estimated uncertainty in 10 millisecond steps
		TimeAuthority  uint `bits:"1"`  // 0 = No Time Authority, 1 = Time Authority
		TAI_UTCDelta   uint `bits:"15"` // Current difference between TAI and UTC in seconds
		TimeZoneOffset uint `bits:"8"`  // The local time zone offset in 15 minute increments
	}

	// Table 5.16: Time Status message parameters
	TimeStatusMessageParameters struct {
		TAISeconds     uint `bits:"40"` // The current TAI time in seconds
		Subsecond      uint `bits:"8"`  // The sub second time in units of 1/256th second (C.1)
		Uncertainty    uint `bits:"8"`  // The estimated uncertainty in 10 millisecond steps (C.1)
		TimeAuthority  uint `bits:"1"`  // 0 = No Time Authority, 1 = Time Authority (C.1)
		TAI_UTCDelta   uint `bits:"15"` // Current difference between TAI and UTC in seconds (C.1)
		TimeZoneOffset uint `bits:"8"`  // The local time zone offset in 15 minute increments (C.1)
	}

	// Table 5.17: Time Zone Set message parameters
	TimeZoneSetMessageParameters struct {
		TimeZoneOffsetNew uint `bits:"8"`  // Upcoming local time zone offset
		TAIOfZoneChange   uint `bits:"40"` // TAI Seconds time of the upcoming Time Zone Offset change
	}

	// Table 5.18: Time Zone Status message parameters
	TimeZoneStatusMessageParameters struct {
		TimeZoneOffsetCurrent uint `bits:"8"`  // Current local time zone offset
		TimeZoneOffsetNew     uint `bits:"8"`  // Upcoming local time zone offset
		TAIOfZoneChange       uint `bits:"40"` // TAI Seconds time of the upcoming Time Zone Offset change
	}

	// Table 5.19: TAI UTC Delta Set message parameters
	TAI_UTCDeltaSetMessageParameters struct {
		TAI_UTCDeltaNew  uint `bits:"15"` // Upcoming difference between TAI and UTC in seconds
		Padding          uint `bits:"1"`  // Always 0b0. Other values are Prohibited.
		TAIOfDeltaChange uint `bits:"40"` // TAI Seconds time of the upcoming TAI UTC Delta change
	}

	// Table 5.20: TAI UTC Delta Status message parameters
	TAI_UTCDeltaStatusMessageParameters struct {
		TAI_UTCDeltaCurrent uint `bits:"15"` // Current difference between TAI and UTC in seconds
		Padding1            uint `bits:"1"`  // Always 0b0. Other values are Prohibited.
		TAI_UTCDeltaNew     uint `bits:"15"` // Upcoming difference between TAI and UTC in seconds
		Padding2            uint `bits:"1"`  // Always 0b0. Other values are Prohibited.
		TAIOfDeltaChange    uint `bits:"40"` // TAI Seconds time of the upcoming TAI UTC Delta change
	}

	// Table 5.21: Time Role Set message parameters
	TimeRoleSetMessageParameters struct {
		TimeRole uint `bits:"8"` // The Time Role for the element
	}

	// Table 5.22: Time Role Status message parameters
	TimeRoleStatusMessageParameters struct {
		TimeRole uint `bits:"8"` // The Time Role for the element
	}

	// Table 5.23: Scene Store message parameters
	SceneStoreMessageParameters struct {
		SceneNumber uint `bits:"16"` // The number of the scene to be stored.
	}

	// Table 5.24: Scene Store Unacknowledged message parameters
	SceneStoreUnacknowledgedMessageParameters struct {
		SceneNumber uint `bits:"16"` // The number of the scene to be stored.
	}

	// Table 5.25: Scene Recall message parameters
	SceneRecallMessageParameters struct {
		SceneNumber    uint                                    `bits:"16"` // The number of the scene to be recalled.
		TID            uint                                    `bits:"8"`  // Transaction Identifier.
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"`  // Message execution delay in 5 millisecond steps. (C.1)
	}

	// Table 5.26: Scene Recall Unacknowledged message parameters
	SceneRecallUnacknowledgedMessageParameters struct {
		SceneNumber    uint                                    `bits:"16"` // The number of the scene to be recalled.
		TID            uint                                    `bits:"8"`  // Transaction Identifier.
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"`  // Message execution delay in 5 millisecond steps. (C.1)
	}

	// Table 5.27: Scene Status message parameters
	SceneStatusMessageParameters struct {
		StatusCode    uint `bits:"8"`  // Defined in 5.2.2.11
		CurrentScene  uint `bits:"16"` // Scene Number of a current scene.
		TargetScene   uint `bits:"16"` // Scene Number of a target scene. (Optional)
		RemainingTime uint `bits:"8"`  // Format as defined in Section 3.1.3. (C.1)
	}

	// Table 5.28: Scene Register Status message parameters
	SceneRegisterStatusMessageParameters struct {
		StatusCode   uint `bits:"8"`        // Defined in Section 5.2.2.11.
		CurrentScene uint `bits:"16"`       // Scene Number of a current scene
		Scenes       uint `bits:"variable"` // A list of scenes stored within an element
	}

	// Table 5.29: Scene Delete message parameter
	SceneDeleteMessageParameter struct {
		SceneNumber uint `bits:"16"` // The number of the scene to be deleted.
	}

	// Table 5.30: Scene Delete Unacknowledged message parameters
	SceneDeleteUnacknowledgedMessageParameters struct {
		SceneNumber uint `bits:"16"` // The number of the scene to be deleted.
	}

	// Table 5.32: Scheduler Status message parameters
	SchedulerStatusMessageParameters struct {
		Schedules uint `bits:"16"` // Bit field indicating defined Actions in the Schedule Register
	}

	// Table 5.33: Scheduler Action Get message parameters
	SchedulerActionGetMessageParameters struct {
		Index uint `bits:"8"` // Index of the Schedule Register entry to get
	}

	// Table 5.34: Scheduler Action Set message parameters
	SchedulerActionSetMessageParameters struct {
		Index            uint `bits:"4"`  // Index of the Schedule Register entry to set
		ScheduleRegister uint `bits:"76"` // Bit field defining an entry in the Schedule Register (see Section 5.1.4.2)
	}

	// Table 5.35: Scheduler Action Set Unacknowledged message parameters
	SchedulerActionSetUnacknowledgedMessageParameters struct {
		Index            uint `bits:"4"`  // Index of the Schedule Register entry to set
		ScheduleRegister uint `bits:"76"` // Bit field defining an entry in the Schedule Register (see Section 5.1.4.2)
	}

	// Table 5.36: Scheduler Action Status message parameters
	SchedulerActionStatusMessageParameters struct {
		Index            uint `bits:"4"`  // Enumerates (selects) a Schedule Register entry
		ScheduleRegister uint `bits:"76"` // Bit field defining an entry in the Schedule Register (see Section 5.1.4.2)
	}

	// Table 6.53: Light Lightness Set message parameters
	LightLightnessSetMessageParameters struct {
		Lightness      uint                                    `bits:"16"` // The target value of the Light Lightness Actual state.
		TID            uint                                    `bits:"8"`  // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"`  // Message execution delay in 5 millisecond steps. (C.1)
	}

	// Table 6.54: Light Lightness Set Unacknowledged message parameters
	LightLightnessSetUnacknowledgedMessageParameters struct {
		Lightness      uint                                    `bits:"16"` // The target value of the Light Lightness Actual state.
		TID            uint                                    `bits:"8"`  // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"`  // Message execution delay in 5 millisecond steps. (C.1)
	}

	// Table 6.55: Light Lightness Status message parameters
	LightLightnessStatusMessageParameters struct {
		PresentLightness uint `bits:"16"`  // The present value of the Light Lightness Actual state.
		TargetLightness  uint `bits:"t16"` // The target value of the Light Lightness Actual state. (Optional)
		RemainingTime    uint `bits:"8"`   // Format as defined in Section 3.1.3. (C.1)
	}

	// Table 6.56: Light Lightness Linear Set message parameters
	LightLightnessLinearSetMessageParameters struct {
		Lightness      uint                                    `bits:"16"` // The target value of the Light Lightness Linear state.
		TID            uint                                    `bits:"8"`  // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"`  // Message execution delay in 5 millisecond steps. (C.1)
	}

	// Table 6.57: Light Lightness Linear Set Unacknowledged message parameters
	LightLightnessLinearSetUnacknowledgedMessageParameters struct {
		Lightness      uint                                    `bits:"16"` // The target value of the Light Lightness Linear state.
		TID            uint                                    `bits:"8"`  // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"`  // Message execution delay in 5 millisecond steps. (C.1)
	}

	// Table 6.58: Light Lightness Linear Status message parameters
	LightLightnessLinearStatusMessageParameters struct {
		PresentLightness uint `bits:"16"`  // The present value of the Light Lightness Linear state
		TargetLightness  uint `bits:"t16"` // The target value of the Light Lightness Linear state (Optional)
		RemainingTime    uint `bits:"8"`   // Format as defined in Section 3.1.3 (C.1)
	}

	// Table 6.59: Light Lightness Last Status message parameters
	LightLightnessLastStatusMessageParameters struct {
		Lightness uint `bits:"16"` // The value of the Light Lightness Last
	}

	// Table 6.60: Light Lightness Default Set message parameters
	LightLightnessDefaultSetMessageParameters struct {
		Lightness uint `bits:"16"` // The value of the Light Lightness Default state
	}

	// Table 6.61: Light Lightness Default Set Unacknowledged message parameters
	LightLightnessDefaultSetUnacknowledgedMessageParameters struct {
		Lightness uint `bits:"16"` // The value of the Light Lightness Default state
	}

	// Table 6.62: Light Lightness Default Status message parameters
	LightLightnessDefaultStatusMessageParameters struct {
		Lightness uint `bits:"16"` // The value of the Light Lightness Default state
	}

	// Table 6.63: Light Lightness Range Set message parameters
	LightLightnessRangeSetMessageParameters struct {
		RangeMin uint `bits:"16"` // The value of the Lightness Range Min field of the Light Lightness Range state
		RangeMax uint `bits:"16"` // The value of the Lightness Range Max field of the Light Lightness Range state
	}

	// Table 6.64: Light Lightness Range Set Unacknowledged message parameters
	LightLightnessRangeSetUnacknowledgedMessageParameters struct {
		RangeMin uint `bits:"16"` // The value of the Lightness Range Min field of the Light Lightness Range state
		RangeMax uint `bits:"16"` // The value of the Lightness Range Max field of the Light Lightness Range state
	}

	// Table 6.65: Light Lightness Range Status message parameters
	LightLightnessRangeStatusMessageParameters struct {
		StatusCode uint `bits:"8"`  // Status Code for the requesting message.
		RangeMin   uint `bits:"16"` // The value of the Lightness Range Min field of the Light Lightness Range state
		RangeMax   uint `bits:"16"` // The value of the Lightness Range Max field of the Light Lightness Range state
	}

	// Table 6.66: Light CTL Set message parameters
	LightCTLSetMessageParameters struct {
		CTLLightness   uint                                    `bits:"16"` // The target value of the Light CTL Lightness state.
		CTLTemperature uint                                    `bits:"16"` // The target value of the Light CTL Temperature state.
		CTLDeltaUV     int                                     `bits:"16"` // The target value of the Light CTL Delta UV state.
		TID            uint                                    `bits:"8"`  // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"`  // Message execution delay in 5 millisecond steps. (C.1)
	}

	// Table 6.67: Light CTL Set Unacknowledged message parameters
	LightCTLSetUnacknowledgedMessageParameters struct {
		CTLLightness   uint                                    `bits:"16"` // The target value of the Light CTL Lightness state.
		CTLTemperature uint                                    `bits:"16"` // The target value of the Light CTL Temperature state.
		CTLDeltaUV     uint                                    `bits:"16"` // The target value of the Light CTL Delta UV state.
		TID            uint                                    `bits:"8"`  // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"`  // Message execution delay in 5 millisecond steps. (C.1)
	}

	// C.1: If the Target CTL Lightness field is present, the Target CTL Temperature and the Remaining Time fields shall also be present; otherwise these fields shall not be present. Table 6.68: Light CTL Status message parameters
	LightCTLStatusMessageParameters struct {
		PresentCTLLightness   uint `bits:"16"`  // The present value of the Light CTL Lightness state
		PresentCTLTemperature uint `bits:"16"`  // The present value of the Light CTL Temperature state
		TargetCTLLightness    uint `bits:"t16"` // The target value of the Light CTL Lightness state (Optional)
		TargetCTLTemperature  uint `bits:"16"`  // The target value of the Light CTL Temperature state (C.1)
		RemainingTime         uint `bits:"8"`   // Format as defined in Section 3.1.3 (C.1)
	}

	// Table 6.69: Light CTL Temperature Set message parameters
	LightCTLTemperatureSetMessageParameters struct {
		CTLTemperature uint                                    `bits:"16"` // The target value of the Light CTL Temperature state.
		CTLDeltaUV     int                                     `bits:"16"` // The target value of the Light CTL Delta UV state.
		TID            uint                                    `bits:"8"`  // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"`  // Message execution delay in 5 millisecond steps. (C.1)
	}

	// Table 6.70: Light CTL Temperature Set Unacknowledged message parameters
	LightCTLTemperatureSetUnacknowledgedMessageParameters struct {
		CTLTemperature uint                                    `bits:"16"` // The target value of the Light CTL Temperature state.
		CTLDeltaUV     int                                     `bits:"16"` // The target value of the Light CTL Delta UV state.
		TID            uint                                    `bits:"8"`  // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"`  // Message execution delay in 5 millisecond steps. (C.1)
	}

	// C.1: If the Target CTL Temperature field is present, the Target CTL Delta UV field and the Remaining Time field shall also be present; otherwise these fields shall not be present. Table 6.71: Light CTL Temperature Status message parameters
	LightCTLTemperatureStatusMessageParameters struct {
		PresentCTLTemperature uint `bits:"16"`  // The present value of the Light CTL Temperature state
		PresentCTLDeltaUV     int  `bits:"16"`  // The present value of the Light CTL Delta UV state
		TargetCTLTemperature  uint `bits:"t16"` // The target value of the Light CTL Temperature state (Optional)
		TargetCTLDeltaUV      uint `bits:"16"`  // The target value of the Light CTL Delta UV state (C.1)
		RemainingTime         uint `bits:"8"`   // Format as defined in Section 3.1.3 (C.1)
	}

	// Table 6.72: Light CTL Temperature Range Set message parameters
	LightCTLTemperatureRangeSetMessageParameters struct {
		RangeMin uint `bits:"16"` // The value of the Temperature Range Min field of the Light CTL Temperature Range state
		RangeMax uint `bits:"16"` // The value of the Temperature Range Max field of the Light CTL Temperature Range state
	}

	// Table 6.73: Light CTL Temperature Range Set Unacknowledged message parameters
	LightCTLTemperatureRangeSetUnacknowledgedMessageParameters struct {
		RangeMin uint `bits:"16"` // The value of the Temperature Range Min field of the Light CTL Temperature Range state
		RangeMax uint `bits:"16"` // The value of the Temperature Range Max field of the Light CTL Temperature Range state
	}

	// Table 6.74: Light CTL Temperature Range Status message parameters
	LightCTLTemperatureRangeStatusMessageParameters struct {
		StatusCode uint `bits:"8"`  // Status Code for the requesting message.
		RangeMin   uint `bits:"16"` // The value of the Temperature Range Min field of the Light CTL Temperature Range state
		RangeMax   uint `bits:"16"` // The value of the Temperature Range Max field of the Light CTL Temperature Range state
	}

	// Table 6.75: Light CTL Default Set message parameters
	LightCTLDefaultSetMessageParameters struct {
		Lightness   uint `bits:"16"` // The value of the Light Lightness Default state
		Temperature uint `bits:"16"` // The value of the Light CTL Temperature Default state
		DeltaUV     int  `bits:"16"` // The value of the Light CTL Delta UV Default state
	}

	// Table 6.76: Light CTL Default Set Unacknowledged message parameters
	LightCTLDefaultSetUnacknowledgedMessageParameters struct {
		Lightness   uint `bits:"16"` // The value of the Light Lightness Default state
		Temperature uint `bits:"16"` // The value of the Light CTL Temperature Default state
		DeltaUV     int  `bits:"16"` // The value of the Light CTL Delta UV Default state
	}

	// Table 6.77: Light CTL Default Status message parameters
	LightCTLDefaultStatusMessageParameters struct {
		Lightness   uint `bits:"16"` // The value of the Light Lightness Default state
		Temperature uint `bits:"16"` // The value of the Light CTL Temperature Default state
		DeltaUV     int  `bits:"16"` // The value of the Light CTL Delta UV Default state
	}

	// Table 6.78: Light HSL Set message parameters
	LightHSLSetMessageParameters struct {
		HSLLightness   uint                                    `bits:"16"` // The target value of the Light HSL Lightness state
		HSLHue         uint                                    `bits:"16"` // The target value of the Light HSL Hue state
		HSLSaturation  uint                                    `bits:"16"` // The target value of the Light HSL Saturation state
		TID            uint                                    `bits:"8"`  // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"`  // Message execution delay in 5 millisecond steps (C.1)
	}

	// Table 6.79: Light HSL Set Unacknowledged message parameters
	LightHSLSetUnacknowledgedMessageParameters struct {
		HSLLightness   uint                                    `bits:"16"` // The target value of the Light HSL Lightness state
		HSLHue         uint                                    `bits:"16"` // The target value of the Light HSL Hue state
		HSLSaturation  uint                                    `bits:"16"` // The target Light HSL Saturation state
		TID            uint                                    `bits:"8"`  // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"`  // Message execution delay in 5 millisecond steps (C.1)
	}

	// Table 6.80: Light HSL Status message parameters
	LightHSLStatusMessageParameters struct {
		HSLLightness  uint `bits:"16"` // The present value of the Light HSL Lightness state
		HSLHue        uint `bits:"16"` // The present value of the Light HSL Hue state
		HSLSaturation uint `bits:"16"` // The present value of the Light HSL Saturation state
		RemainingTime uint `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
	}

	// Table 6.81: Light HSL Target Status message parameters
	LightHSLTargetStatusMessageParameters struct {
		HSLLightnessTarget  uint `bits:"16"` // The target value of the Light HSL Lightness state
		HSLHueTarget        uint `bits:"16"` // The target value of the Light HSL Hue state
		HSLSaturationTarget uint `bits:"16"` // The target Light HSL Saturation state
		RemainingTime       uint `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
	}

	// Table 6.82: Light HSL Hue Set message parameters
	LightHSLHueSetMessageParameters struct {
		Hue            uint                                    `bits:"16"` // The target value of the Light HSL Hue state.
		TID            uint                                    `bits:"8"`  // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"`  // Message execution delay in 5 millisecond steps. (C.1)
	}

	// Table 6.83: Light HSL Hue Set Unacknowledged message parameters
	LightHSLHueSetUnacknowledgedMessageParameters struct {
		Hue            uint                                    `bits:"16"` // The target value of the Light HSL Hue state.
		TID            uint                                    `bits:"8"`  // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"`  // Message execution delay in 5 millisecond steps. (C.1)
	}

	// Table 6.84: Light HSL Hue Status message parameters
	LightHSLHueStatusMessageParameters struct {
		PresentHue    uint `bits:"16"` // The present value of the Light HSL Hue state
		TargetHue     uint `bits:"16"` // The target value of the Light HSL Hue state (Optional)
		RemainingTime uint `bits:"8"`  // Format as defined in Section 3.1.3 (C.1)
	}

	// Table 6.85: Light HSL Saturation Set message parameters
	LightHSLSaturationSetMessageParameters struct {
		Saturation     uint                                    `bits:"16"` // The target value of the Light HSL Saturation state.
		TID            uint                                    `bits:"8"`  // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"`  // Message execution delay in 5 millisecond steps. (C.1)
	}

	// Table 6.86: Light HSL Saturation Set Unacknowledged message parameters
	LightHSLSaturationSetUnacknowledgedMessageParameters struct {
		Saturation     uint                                    `bits:"16"` // The target value of the Light HSL Saturation state.
		TID            uint                                    `bits:"8"`  // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"`  // Message execution delay in 5 millisecond steps. (C.1)
	}

	// Table 6.87: Light HSL Saturation Status message parameters
	LightHSLSaturationStatusMessageParameters struct {
		PresentSaturation uint `bits:"16"` // The present value of the Light HSL Saturation state.
		TargetSaturation  uint `bits:"16"` // The target value of the Light HSL Saturation state. (Optional)
		RemainingTime     uint `bits:"8"`  // Format as defined in Section 3.1.3. (C.1)
	}

	// Table 6.88: Light HSL Default Set message parameters
	LightHSLDefaultSetMessageParameters struct {
		Lightness  uint `bits:"16"` // The value of the Light Lightness Default state
		Hue        uint `bits:"16"` // The value of the Light HSL Hue Default state
		Saturation uint `bits:"16"` // The value of the Light HSL Saturation Default state
	}

	// Table 6.89: Light HSL Default Set Unacknowledged message parameters The Lightness field identifies the Light Lightness Default state of the element.
	LightHSLDefaultSetUnacknowledgedMessageParameters struct {
		Lightness  uint `bits:"16"` // The value of the Light Lightness Default state
		Hue        uint `bits:"16"` // The value of the Light HSL Hue Default state
		Saturation uint `bits:"16"` // The value of the Light HSL Saturation Default state
	}

	// Table 6.90: Light HSL Default Status message parameters
	LightHSLDefaultStatusMessageParameters struct {
		Lightness  uint `bits:"16"` // The value of the Light Lightness Default state
		Hue        uint `bits:"16"` // The value of the Light HSL Hue Default state
		Saturation uint `bits:"16"` // The value of the Light HSL Saturation Default state
	}

	// Table 6.91: Light HSL Range Set message parameters
	LightHSLRangeSetMessageParameters struct {
		HueRangeMin        uint `bits:"16"` // The value of the Hue Range Min field of the Light HSL Hue Range state
		HueRangeMax        uint `bits:"16"` // The value of the Hue Range Max field of the Light HSL Hue Range state
		SaturationRangeMin uint `bits:"16"` // The value of the Saturation Range Min field of the Light HSL Saturation Range state
		SaturationRangeMax uint `bits:"16"` // The value of the Saturation Range Max field of the Light HSL Saturation Range state
	}

	// Table 6.92: Light HSL Range Set Unacknowledged message parameters
	LightHSLRangeSetUnacknowledgedMessageParameters struct {
		HueRangeMin        uint `bits:"16"` // The value of the Hue Range Min field of the Light HSL Hue Range state
		HueRangeMax        uint `bits:"16"` // The value of the Hue Range Max field of the Light HSL Hue Range state
		SaturationRangeMin uint `bits:"16"` // The value of the Saturation Range Min field of the Light HSL Saturation Range state
		SaturationRangeMax uint `bits:"16"` // The value of the Saturation Range Max field of the Light HSL Saturation Range state
	}

	// Table 6.93: Light HSL Range Status message parameters
	LightHSLRangeStatusMessageParameters struct {
		StatusCode         uint `bits:"8"`  // Status Code for the requesting message.
		HueRangeMin        uint `bits:"16"` // The value of the Hue Range Min field of the Light HSL Hue Range state
		HueRangeMax        uint `bits:"16"` // The value of the Hue Range Max field of the Light HSL Hue Range state
		SaturationRangeMin uint `bits:"16"` // The value of the Saturation Range Min field of the Light HSL Saturation Range state
		SaturationRangeMax uint `bits:"16"` // The value of the Saturation Range Max field of the Light HSL Saturation Range state
	}

	// Table 6.94: Light xyL Set message parameters
	LightXyLSetMessageParameters struct {
		XyLLightness   uint                                    `bits:"16"` // The target value of the Light xyL Lightness state
		XyLX           uint                                    `bits:"16"` // The target value of the Light xyL x state
		XyLY           uint                                    `bits:"16"` // The target value of the Light xyL y state
		TID            uint                                    `bits:"8"`  // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"`  // Message execution delay in 5 millisecond steps (C.1)
	}

	// Table 6.95: Light xyL Set Unacknowledged message parameters
	LightXyLSetUnacknowledgedMessageParameters struct {
		XyLLightness   uint                                    `bits:"16"` // The target value of the Light xyL Lightness state
		XyLX           uint                                    `bits:"16"` // The target value of the Light xyL x state
		XyLY           uint                                    `bits:"16"` // The target value of the Light xyL y state
		TID            uint                                    `bits:"8"`  // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"`  // Message execution delay in 5 millisecond steps (C.1)
	}

	// Table 6.96: Light xyL Status Unacknowledged message parameters
	LightXyLStatusUnacknowledgedMessageParameters struct {
		XyLLightness  uint `bits:"16"` // The present value of the Light xyL Lightness state
		XyLX          uint `bits:"16"` // The present value of the Light xyL x state
		XyLY          uint `bits:"16"` // The present value of the Light xyL y state
		RemainingTime uint `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
	}

	// Table 6.97: Light xyL Target Status Unacknowledged message parameters
	LightXyLTargetStatusUnacknowledgedMessageParameters struct {
		TargetXyLLightness uint `bits:"16"` // The target value of the Light xyL Lightness state
		TargetXyLX         uint `bits:"16"` // The target value of the Light xyL x state
		TargetXyLY         uint `bits:"16"` // The target value of the Light xyL y state
		RemainingTime      uint `bits:"8"`  // Format as defined in Section 3.1.3. (Optional)
	}

	// Table 6.98: Light HSL Default Set message parameters
	LightXyLDefaultSetMessageParameters struct {
		Lightness uint `bits:"16"` // The value of the Light Lightness Default state
		XyLX      uint `bits:"16"` // The value of the Light xyL x Default state
		XyLY      uint `bits:"16"` // The value of the Light xyL y Default state
	}

	// Table 6.99: Light xyL Default Set Unacknowledged message parameters The Lightness field identifies the Light Lightness Default state of the element.
	LightXyLDefaultSetUnacknowledgedMessageParameters struct {
		Lightness uint `bits:"16"` // The value of the Light Lightness Default state
		XyLX      uint `bits:"16"` // The value of the Light xyL x Default state
		XyLY      uint `bits:"16"` // The value of the Light xyL y Default state
	}

	// Table 6.100: Light xyL Default Status message parameters
	LightXyLDefaultStatusMessageParameters struct {
		Lightness uint `bits:"16"` // The value of the Light Lightness Default state
		XyLX      uint `bits:"16"` // The value of the Light xyL x Default state
		XyLY      uint `bits:"16"` // The value of the Light xyL y Default state
	}

	// Table 6.101: Light xyL Range Set message parameters
	LightXyLRangeSetMessageParameters struct {
		XyLXRangeMin uint `bits:"16"` // The value of the xyL x Range Min field of the Light xyL x Range state
		XyLXRangeMax uint `bits:"16"` // The value of the xyL x Range Max field of the Light xyL x Range state
		XyLYRangeMin uint `bits:"16"` // The value of the xyL y Range Min field of the Light xyL y Range state
		XyLYRangeMax uint `bits:"16"` // The value of the xyL y Range Max field of the Light xyL y Range state
	}

	// Table 6.102: Light xyL Range Set Unacknowledged message parameters
	LightXyLRangeSetUnacknowledgedMessageParameters struct {
		XyLXRangeMin uint `bits:"16"` // The value of the xyL x Range Min field of the Light xyL x Range state
		XyLXRangeMax uint `bits:"16"` // The value of the xyL x Range Max field of the Light xyL x Range state
		XyLYRangeMin uint `bits:"16"` // The value of the xyL y Range Min field of the Light xyL y Range state
		XyLYRangeMax uint `bits:"16"` // The value of the xyL y Range Max field of the Light xyL y Range state
	}

	// Table 6.103: Light xyL Range Status message parameters
	LightXyLRangeStatusMessageParameters struct {
		StatusCode   uint `bits:"8"`  // Status Code for the requesting message.
		XyLXRangeMin uint `bits:"16"` // The value of the xyL x Range Min field of the Light xyL x Range state
		XyLXRangeMax uint `bits:"16"` // The value of the xyL x Range Max field of the Light xyL x Range state
		XyLYRangeMin uint `bits:"16"` // The value of the xyL y Range Min field of the Light xyL y Range state
		XyLYRangeMax uint `bits:"16"` // The value of the xyL y Range Max field of the Light xyL y Range state
	}

	// Table 6.104: Light LC Mode Set message parameters
	LightLCModeSetMessageParameters struct {
		Mode uint `bits:"8"` // The target value of the Light LC Mode state
	}

	// Table 6.105: Light LC Mode Set Unacknowledged message parameters
	LightLCModeSetUnacknowledgedMessageParameters struct {
		Mode uint `bits:"8"` // The target value of the Light LC Mode state
	}

	// Table 6.106: Light LC Mode Status message parameters
	LightLCModeStatusMessageParameters struct {
		Mode uint `bits:"8"` // The present value of the Light LC Mode state
	}

	// Table 6.107: Light LC OM Set message parameters
	LightLCOMSetMessageParameters struct {
		Mode uint `bits:"8"` // The target value of the Light LC Occupancy Mode state
	}

	// Table 6.108: Light LC OM Set Unacknowledged message parameters
	LightLCOMSetUnacknowledgedMessageParameters struct {
		Mode uint `bits:"8"` // The target value of the Light LC Occupancy Mode state
	}

	// Table 6.109: Light LC OM Status message parameters
	LightLCOMStatusMessageParameters struct {
		Mode uint `bits:"8"` // The present value of the Light LC Occupancy Mode state
	}

	// Table 6.110: Light LC Light OnOff Set message parameters
	LightLCLightOnOffSetMessageParameters struct {
		LightOnOff     uint                                    `bits:"8"` // The target value of the Light LC Light OnOff state
		TID            uint                                    `bits:"8"` // Transaction Identifier
		TransitionTime GenericDefaultTransitionTimeStateFormat `bits:"8"` // Format as defined in Section 3.1.3. (Optional)
		Delay          uint                                    `bits:"8"` // Message execution delay in 5 millisecond steps (C.1)
	}

	// Table 6.112: Light LC Light OnOff Status message parameters
	LightLCLightOnOffStatusMessageParameters struct {
		PresentLightOnOff uint `bits:"8"` // The present value of the Light LC Light OnOff state
		TargetLightOnOff  uint `bits:"8"` // The target value of the Light LC Light OnOff state (Optional)
		RemainingTime     uint `bits:"8"` // Format as defined in Section 3.1.3. (C.1)
	}

	// Table 6.113: Light LC Property Get message parameters
	LightLCPropertyGetMessageParameters struct {
		LightLCPropertyID uint `bits:"16"` // Property ID identifying a Light LC Property.
	}

	// Table 6.114: Light LC Property Set message parameters
	LightLCPropertySetMessageParameters struct {
		LightLCPropertyID    uint `bits:"16"`       // Property ID identifying a Light LC Property.
		LightLCPropertyValue uint `bits:"variable"` // Raw value for the Light LC Property
	}

	// Table 6.115: Light LC Property Set Unacknowledged message parameters
	LightLCPropertySetUnacknowledgedMessageParameters struct {
		LightLCPropertyID    uint `bits:"16"`       // Property ID identifying a Light LC Property.
		LightLCPropertyValue uint `bits:"variable"` // Raw value for the Light LC Property
	}

	// Table 6.116: Light LC Property Status message parameters
	LightLCPropertyStatusMessageParameters struct {
		LightLCPropertyID    uint `bits:"16"`       // Property ID identifying a Light LC Property.
		LightLCPropertyValue uint `bits:"variable"` // Raw value for the Light LC Property
	}
)
