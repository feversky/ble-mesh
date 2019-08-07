package def

type (
	Range struct {
		Start uint
		End   uint
	}

	// Table 3.7: Network PDU field definitions
	NetworkPDUFieldDefinitions struct {
		IVI          uint   `bits:"1"`     // Least significant bit of IV Index
		NID          uint   `bits:"7"`     // Value derived from the NetKey used to identify the Encryption Key and Privacy Key used to secure this PDU
		CTL          uint   `bits:"1"`     // Network Control
		TTL          uint   `bits:"7"`     // Time To Live
		SEQ          uint   `bits:"24"`    // Sequence Number
		SRC          uint   `bits:"16"`    // Source Address
		DST          uint   `bits:"16"`    // Destination Address
		TransportPDU []byte `bits:"8"`     // Transport Protocol Data Unit
		NetMIC       uint   `bits:"32/64"` // Message Integrity Check for Network
	}

	// Table 3.10: Unsegmented Access message format
	UnsegmentedAccessMessageFormat struct {
		SEG                     uint   `bits:"1"` // 0 = Unsegmented Message
		AKF                     uint   `bits:"1"` // Application Key Flag
		AID                     uint   `bits:"6"` // Application key identifier
		UpperTransportAccessPDU []uint `bits:"8"` // 40 to 120, The Upper Transport Access PDU
	}

	// Table 3.11: Segmented Access message format
	SegmentedAccessMessageFormat struct {
		SEG      uint   `bits:"1"`  // 1 = Segmented Message
		AKF      uint   `bits:"1"`  // Application Key Flag
		AID      uint   `bits:"6"`  // Application key identifier
		SZMIC    uint   `bits:"1"`  // Size of TransMIC
		SeqZero  uint   `bits:"13"` // Least significant bits of SeqAuth
		SegO     uint   `bits:"5"`  // Segment Offset number
		SegN     uint   `bits:"5"`  // Last Segment number
		SegmentM []uint `bits:"8"`  // 8 to 96, Segment m of the Upper Transport Access PDU
	}

	// Table 3.12: Unsegmented Control message format
	UnsegmentedControlMessageFormat struct {
		SEG    uint `bits:"1"` // 0 = Unsegmented Message
		Opcode uint `bits:"7"` // 0x00 = Segment Acknowledgment 0x01 to 0x7F = Opcode of the Transport Control message
	}

	// Table 3.13: Segment Acknowledgment message format
	SegmentAcknowledgmentMessageFormat struct {
		SEG      uint `bits:"1"`  // 0 = Unsegmented Message
		Opcode   uint `bits:"7"`  // 0x00 = Segment Acknowledgment Message
		OBO      uint `bits:"1"`  // Friend on behalf of a Low Power node
		SeqZero  uint `bits:"13"` // SeqZero of the Upper Transport PDU
		RFU      uint `bits:"2"`  // Reserved for Future Use
		BlockAck uint `bits:"32"` // Block acknowledgment for segments
	}

	// Table 3.14: Segmented Control message format
	SegmentedControlMessageFormat struct {
		SEG      uint   `bits:"1"`  // 1 = Segmented Message
		Opcode   uint   `bits:"7"`  // 0x00 = Reserved 0x01 to 0x7F = Opcode of the Transport Control message
		RFU      uint   `bits:"1"`  // Reserved for Future Use
		SeqZero  uint   `bits:"13"` // Least significant bits of SeqAuth
		SegO     uint   `bits:"5"`  // Segment Offset number
		SegN     uint   `bits:"5"`  // Last Segment number
		SegmentM []uint `bits:"8"`  // 8 to 64, Segment m of the Upper Transport Control PDU
	}

	// Table 3.17: Friend Poll message parameters
	FriendPollMessageParameters struct {
		Padding uint `bits:"07"` // 0b0000000. All other values are Prohibited.
		FSN     uint `bits:"1"`  // Friend Sequence Number, used to acknowledge receipt of previous messages from the Friend node to the Low Power node
	}

	// Table 3.18: Friend Update message parameters
	FriendUpdateMessageParameters struct {
		Flags   uint `bits:"8"`  // Contains the IV Update Flag and the Key Refresh Flag
		IVIndex uint `bits:"32"` // The current IV Index value known by the Friend node
		MD      uint `bits:"8"`  // Indicates if the Friend Queue is empty or not.
	}

	// Table 3.21: Friend Request message parameters
	FriendRequestMessageParameters struct {
		Criteria        uint `bits:"8"`  // The criteria that a Friend node should support in order to participate in friendship negotiation
		ReceiveDelay    uint `bits:"8"`  // Receive delay requested by the Low Power node
		PollTimeout     uint `bits:"24"` // The initial value of the PollTimeout timer set by the Low Power node
		PreviousAddress uint `bits:"16"` // Unicast address of the primary element of the previous friend
		NumElements     uint `bits:"8"`  // Number of elements in the Low Power node
		LPNCounter      uint `bits:"16"` // Number of Friend Request messages that the Low Power node has sent
	}

	// Table 3.22: Criteria field format
	CriteriaFieldFormat struct {
		RFU                 uint `bits:"1"` // Reserved for Future Use
		RSSIFactor          uint `bits:"2"` // The contribution of the RSSI measured by the Friend node used in Friend Offer Delay calculations
		ReceiveWindowFactor uint `bits:"2"` // The contribution of the supported Receive Window used in Friend Offer Delay calculations
		MinQueueSizeLog     uint `bits:"3"` // The minimum number of messages that the Friend node can store in its Friend Queue
	}

	// Table 3.29: Friend Offer parameters
	FriendOfferParameters struct {
		ReceiveWindow        uint `bits:"8"`  // Receive Window value supported by the Friend node
		QueueSize            uint `bits:"8"`  // Queue Size available on the Friend node
		SubscriptionListSize uint `bits:"8"`  // Size of the Subscription List that can be supported by a Friend node for a Low Power node
		RSSI                 uint `bits:"8"`  // RSSI measured by the Friend node
		FriendCounter        uint `bits:"16"` // Number of Friend Offer messages that the Friend node has sent
	}

	// Table 3.31: Friend Clear parameters
	FriendClearParameters struct {
		LPNAddress uint `bits:"16"` // The unicast address of the Low Power node being removed
		LPNCounter uint `bits:"16"` // Value of the LPNCounter of new relationship
	}

	// Table 3.32: Friend Clear Confirm parameters
	FriendClearConfirmParameters struct {
		LPNAddress uint `bits:"16"` // The unicast address of the Low Power node being removed
		LPNCounter uint `bits:"16"` // The value of the LPNCounter of corresponding Friend Clear message
	}

	// Table 3.33: Friend Subscription List Add parameters
	FriendSubscriptionListAddParameters struct {
		TransactionNumber uint   `bits:"8"`  // The number for identifying a transaction
		AddressList       []uint `bits:"16"` // List of group addresses and virtual addresses where N is the number of group addresses and virtual addresses in this message
	}

	// Table 3.34: Friend Subscription List Remove parameters
	FriendSubscriptionListRemoveParameters struct {
		TransactionNumber uint   `bits:"8"`  // The number for identifying a transaction
		AddressList       []uint `bits:"16"` // List of group addresses and virtual addresses where N is the number of group addresses and virtual addresses in this message
	}

	// Table 3.35: Friend Subscription List Confirm parameters
	FriendSubscriptionListConfirmParameters struct {
		TransactionNumber uint `bits:"8"` // The number for identifying a transaction
	}

	// Table 3.36: Heartbeat parameters
	HeartbeatParameters struct {
		RFU      uint `bits:"1"`  // Reserved for Future Use
		InitTTL  uint `bits:"7"`  // Initial TTL used when sending the message
		Features uint `bits:"16"` // Bit field of currently active features of the node
	}

	// Table 3.40: Vendor Model ID format
	VendorModelIDFormat struct {
		bitCompanyIdentifier             uint `bits:"16"` // See [6]
		bitVendorAssignedModelIdentifier uint `bits:"16"` //
	}

	// Table 3.41: Access payload fields
	// AccessPayloadFields struct {
	// 	Opcode uint `bits:"8/16/24"` // Operation Code
	// }

	// Table 3.45: Network nonce format
	NetworkNonceFormat struct {
		NonceType uint `bits:"8"`  // 0x00
		CTLAndTTL uint `bits:"8"`  // See Table 3.46
		SEQ       uint `bits:"24"` // Sequence Number
		SRC       uint `bits:"16"` // Source Address
		Pad       uint `bits:"16"` // 0x0000
		IVIndex   uint `bits:"32"` // IV Index
	}

	// Table 3.46: CTL and TTL field format
	CTLAndTTLFieldFormat struct {
		CTL uint `bits:"1"` // See Section 3.4.4.3
		TTL uint `bits:"7"` // See Section 3.4.4.4
	}

	// Table 3.47: Application nonce format
	ApplicationNonceFormat struct {
		NonceType    uint `bits:"8"`  // 0x01
		ASZMICAndPad uint `bits:"8"`  // See Table 3.48
		SEQ          uint `bits:"24"` // Sequence Number of the Access message (24 lowest bits of SeqAuth in the context of segmented messages)
		SRC          uint `bits:"16"` // Source Address
		DST          uint `bits:"16"` // Destination Address
		IVIndex      uint `bits:"32"` // IV Index
	}

	// Table 3.48: ASZMIC and Pad field format
	ASZMICAndPadFieldFormat struct {
		ASZMIC uint `bits:"1"` // SZMIC if a Segmented Access message or 0 for all other message formats
		Pad    uint `bits:"7"` // 0b0000000
	}

	// Table 3.49: Device nonce format
	DeviceNonceFormat struct {
		NonceType    uint `bits:"8"`  // 0x02
		ASZMICAndPad uint `bits:"8"`  // See Table 3.50
		SEQ          uint `bits:"24"` // Sequence Number of the Access message (24 lowest bits of SeqAuth in the context of segmented messages)
		SRC          uint `bits:"16"` // Source Address
		DST          uint `bits:"16"` // Destination Address
		IVIndex      uint `bits:"32"` // IV Index
	}

	// Table 3.51: Proxy nonce format
	ProxyNonceFormat struct {
		NonceType uint `bits:"8"`  // 0x03
		Pad       uint `bits:"8"`  // 0x00
		SEQ       uint `bits:"24"` // Sequence Number
		SRC       uint `bits:"16"` // Source Address
		Pad1      uint `bits:"16"` // 0x0000
		IVIndex   uint `bits:"32"` // IV Index
	}

	// Table 3.53: Unprovisioned Device beacon format
	UnprovisionedDeviceBeaconFormat struct {
		BeaconType     uint `bits:"8"`   // Unprovisioned Device beacon type (0x00)
		DeviceUUID     uint `bits:"128"` // Device UUID uniquely identifying this device (see Section 3.10.3)
		OOBInformation uint `bits:"16"`  // See Table 3.54
		URIHash        uint `bits:"32"`  // Hash of the associated URI advertised with the URI AD Type (optional field)
	}

	// Table 3.55: Secure Network beacon format:
	SecureNetworkBeaconFormat struct {
		BeaconType          uint `bits:"8"`  // Secure Network beacon (0x01)
		Flags               uint `bits:"8"`  // Contains the Key Refresh Flag and IV Update Flag
		NetworkID           uint `bits:"64"` // Contains the value of the Network ID
		IVIndex             uint `bits:"32"` // Contains the current IV Index
		AuthenticationValue uint `bits:"64"` // Authenticates security network beacon
	}

	// Table 4.2: Composition Data Page 0 fields
	CompositionDataPage0Fields struct {
		CID      uint   `bits:"16"` // Contains a 16-bit company identifier assigned by the Bluetooth SIG (the list is available at [6])
		PID      uint   `bits:"16"` // Contains a 16-bit vendor-assigned product identifier
		VID      uint   `bits:"16"` // Contains a 16-bit vendor-assigned product version identifier
		CRPL     uint   `bits:"16"` // Contains a 16-bit value representing the minimum number of replay protection list entries in a device (see Section 3.8.8)
		Features uint   `bits:"16"` // Contains a bit field indicating the device features, as defined in Table 4.3
		Elements []uint `bits:"8"`  // Contains a sequence of element descriptions
	}

	// Table 4.4: Element description format
	ElementDescriptionFormat struct {
		Loc          uint   `bits:"16"` // Contains a location descriptor
		NumS         uint   `bits:"8"`  // Contains a count of SIG Model IDs in this element
		NumV         uint   `bits:"8"`  // Contains a count of Vendor Model IDs in this element
		SIGModels    []uint `bits:"16"` // Contains a sequence of NumS SIG Model IDs
		VendorModels []uint `bits:"32"` // Contains a sequence of NumV Vendor Model IDs
	}

	// Table 4.5: Publish Period format
	PublishPeriodFormat struct {
		NumberOfSteps  uint `bits:"6" vt:"StepResolutionValues"` // The number of steps
		StepResolution uint `bits:"2" vt:"NumberOfStepsValues"`  // The resolution of the Number of Steps field
	}

	// Table 4.19: Current Fault format
	CurrentFaultFormat struct {
		TestID     uint `bits:"8"` // Identifier of a most recently performed self-test
		FaultArray uint `bits:"N"` // Array of current faults
	}

	// Table 4.22: Registered Fault format
	RegisteredFaultFormat struct {
		TestID     uint `bits:"8"` // Identifier of a most recently performed self-test
		FaultArray uint `bits:"N"` // Array of registered faults
	}

	// Table 4.33: Config Beacon Set message parameters
	ConfigBeaconSetMessageParameters struct {
		Beacon uint `bits:"8"` // New Secure Network Beacon state
	}

	// Table 4.34: Config Beacon Status message parameters
	ConfigBeaconStatusMessageParameters struct {
		Beacon uint `bits:"8" vt:"SecureNetworkBeaconValues"` // Secure Network Beacon state
	}

	// Table 4.35: Config Composition Data Get message parameters
	ConfigCompositionDataGetMessageParameters struct {
		Page uint `bits:"8"` // Page number of the Composition Data
	}

	// Table 4.36: Config Composition Data Status message parameters
	ConfigCompositionDataStatusMessageParameters struct {
		Page uint   `bits:"8"` // Page number of the Composition Data
		Data []uint `bits:"8"` // Composition Data for the identified page
	}

	// Table 4.37: Config Default TTL Set message parameters
	ConfigDefaultTTLSetMessageParameters struct {
		TTL uint `bits:"8"` // New Default TTL value
	}

	// Table 4.38: Config Default TTL Status message parameters
	ConfigDefaultTTLStatusMessageParameters struct {
		TTL uint `bits:"8" vt:"DefaultTTLValues"` // Default TTL
	}

	// Table 4.39: Config GATT Proxy Set message parameters
	ConfigGATTProxySetMessageParameters struct {
		GATTProxy uint `bits:"8"` // New GATT Proxy state
	}

	// Table 4.40: Config GATT Proxy Status message parameters
	ConfigGATTProxyStatusMessageParameters struct {
		GATTProxy uint `bits:"8" vt:"GATTProxyValues"` // GATT Proxy state
	}

	// Table 4.41: Config Relay Set message parameters
	ConfigRelaySetMessageParameters struct {
		Relay                        uint `bits:"8"` // Relay
		RelayRetransmitCount         uint `bits:"3"` // Number of retransmissions on advertising bearer for each Network PDU relayed by the node
		RelayRetransmitIntervalSteps uint `bits:"5"` // Number of 10-millisecond steps between retransmissions
	}

	// Table 4.42: Config Relay Status message parameters
	ConfigRelayStatusMessageParameters struct {
		Relay                        uint `bits:"8" vt:"RelayValues"` // Relay
		RelayRetransmitCount         uint `bits:"3"`                  // Number of retransmissions on advertising bearer for each Network PDU relayed by the node
		RelayRetransmitIntervalSteps uint `bits:"5"`                  // Number of 10-millisecond steps between retransmissions
	}

	// Table 4.43: Config Model Publication Get message parameters
	ConfigModelPublicationGetMessageParameters struct {
		ElementAddress  uint `bits:"16"`    // Address of the element
		ModelIdentifier uint `bits:"16/32"` // SIG Model ID or Vendor Model ID
	}

	// Table 4.44: Config Model Publication Set message parameters
	ConfigModelPublicationSetMessageParameters struct {
		ElementAddress                 uint                `bits:"16"`    // Address of the element
		PublishAddress                 uint                `bits:"16"`    // Value of the publish address
		AppKeyIndex                    uint                `bits:"12"`    // Index of the application key
		CredentialFlag                 uint                `bits:"1"`     // Value of the Friendship Credential Flag
		RFU                            uint                `bits:"3"`     // Reserved for Future Use
		PublishTTL                     uint                `bits:"8"`     // Default TTL value for the outgoing messages
		PublishPeriod                  PublishPeriodFormat `bits:"8"`     // Period for periodic status publishing
		PublishRetransmitCount         uint                `bits:"3"`     // Number of retransmissions for each published message
		PublishRetransmitIntervalSteps uint                `bits:"5"`     // Number of 50-millisecond steps between retransmissions
		ModelIdentifier                uint                `bits:"16/32"` // SIG Model ID or Vendor Model ID
	}

	// Table 4.45: Config Model Publication Virtual Address Set message parameters
	ConfigModelPublicationVirtualAddressSetMessageParameters struct {
		ElementAddress                 uint                `bits:"16"`    // Address of the element
		PublishAddress                 uint                `bits:"128"`   // Value of the Label UUID publish address
		AppKeyIndex                    uint                `bits:"12"`    // Index of the application key
		CredentialFlag                 uint                `bits:"1"`     // Value of the Friendship Credential Flag
		RFU                            uint                `bits:"3"`     // Reserved for Future Use
		PublishTTL                     uint                `bits:"8"`     // Default TTL value for the outgoing messages
		PublishPeriod                  PublishPeriodFormat `bits:"8"`     // Period for periodic status publishing
		PublishRetransmitCount         uint                `bits:"3"`     // Number of retransmissions for each published message
		PublishRetransmitIntervalSteps uint                `bits:"5"`     // Number of 50-millisecond steps between retransmissions
		ModelIdentifier                uint                `bits:"16/32"` // SIG Model ID or Vendor Model ID
	}

	// Table 4.46: Config Model Publication Status message parameters
	ConfigModelPublicationStatusMessageParameters struct {
		Status                         uint                `bits:"8" vt:"SummaryOfStatusCodes"`                  // Status Code for the requesting message
		ElementAddress                 uint                `bits:"16"`                                           // Address of the element
		PublishAddress                 uint                `bits:"16"`                                           // Value of the publish address
		AppKeyIndex                    uint                `bits:"12"`                                           // Index of the application key
		CredentialFlag                 uint                `bits:"1" vt:"PublishFriendshipCredentialFlagValues"` // Value of the Friendship Credential Flag
		RFU                            uint                `bits:"3"`                                            // Reserved for Future Use
		PublishTTL                     uint                `bits:"8" vt:"DefaultTTLValues"`                      // Default TTL value for the outgoing messages
		PublishPeriod                  PublishPeriodFormat `bits:"8"`                                            // Period for periodic status publishing
		PublishRetransmitCount         uint                `bits:"3"`                                            // Number of retransmissions for each published message
		PublishRetransmitIntervalSteps uint                `bits:"5"`                                            // Number of 50-millisecond steps between retransmissions
		ModelIdentifier                uint                `bits:"16/32"`                                        // SIG Model ID or Vendor Model ID
	}

	// Table 4.47: Config Model Subscription Add message parameters
	ConfigModelSubscriptionAddMessageParameters struct {
		ElementAddress  uint `bits:"16"`    // Address of the element
		Address         uint `bits:"16"`    // Value of the address
		ModelIdentifier uint `bits:"16/32"` // SIG Model ID or Vendor Model ID
	}

	// Table 4.48: Config Model Subscription Virtual Address Add message parameters
	ConfigModelSubscriptionVirtualAddressAddMessageParameters struct {
		ElementAddress  uint `bits:"16"`    // Address of the element
		Label           uint `bits:"128"`   // Value of the Label UUID
		ModelIdentifier uint `bits:"16/32"` // SIG Model ID or Vendor Model ID
	}

	// Table 4.49: Config Model Subscription Delete message parameters
	ConfigModelSubscriptionDeleteMessageParameters struct {
		ElementAddress  uint `bits:"16"`    // Address of the element
		Address         uint `bits:"16"`    // Value of the Address
		ModelIdentifier uint `bits:"16/32"` // SIG Model ID or Vendor Model ID
	}

	// Table 4.50: Config Model Subscription Virtual Address Delete message parameters
	ConfigModelSubscriptionVirtualAddressDeleteMessageParameters struct {
		ElementAddress  uint `bits:"16"`    // Address of the element
		Address         uint `bits:"128"`   // Value of the Label UUID
		ModelIdentifier uint `bits:"16/32"` // SIG Model ID or Vendor Model ID
	}

	// Table 4.51: Config Model Subscription Overwrite message parameters
	ConfigModelSubscriptionOverwriteMessageParameters struct {
		ElementAddress  uint `bits:"16"`    // Address of the element
		Address         uint `bits:"16"`    // Value of the Address
		ModelIdentifier uint `bits:"16/32"` // SIG Model ID or Vendor Model ID
	}

	// Table 4.52: Config Model Subscription Virtual Address Overwrite message parameters
	ConfigModelSubscriptionVirtualAddressOverwriteMessageParameters struct {
		ElementAddress  uint `bits:"16"`    // Address of the element
		Address         uint `bits:"128"`   // Value of the Label UUID
		ModelIdentifier uint `bits:"16/32"` // SIG Model ID or Vendor Model ID
	}

	// Table 4.53: Config Model Subscription Delete All message parameters
	ConfigModelSubscriptionDeleteAllMessageParameters struct {
		ElementAddress  uint `bits:"16"`    // Address of the element
		ModelIdentifier uint `bits:"16/32"` // SIG Model ID or Vendor Model ID
	}

	// Table 4.54: Config Model Subscription Status message parameters
	ConfigModelSubscriptionStatusMessageParameters struct {
		Status          uint `bits:"8" vt:"SummaryOfStatusCodes"` // Status Code for the requesting message.
		ElementAddress  uint `bits:"16"`                          // Address of the element
		Address         uint `bits:"16"`                          // Value of the address
		ModelIdentifier uint `bits:"16/32"`                       // SIG Model ID or Vendor Model ID
	}

	// Table 4.55: Config SIG Model Subscription Get message parameters
	ConfigSIGModelSubscriptionGetMessageParameters struct {
		ElementAddress  uint `bits:"16"` // Address of the element
		ModelIdentifier uint `bits:"16"` // SIG Model ID
	}

	// Table 4.56: Config SIG Model Subscription List message parameters
	ConfigSIGModelSubscriptionListMessageParameters struct {
		Status          uint   `bits:"8"`  // Status Code for the requesting message
		ElementAddress  uint   `bits:"16"` // Address of the element
		ModelIdentifier uint   `bits:"16"` // SIG Model ID
		Addresses       []uint `bits:"16"` // A block of all addresses from the Subscription List
	}

	// Table 4.57: Config Vendor Model Subscription Get message parameters
	ConfigVendorModelSubscriptionGetMessageParameters struct {
		ElementAddress  uint `bits:"16"` // Address of the element
		ModelIdentifier uint `bits:"32"` // Vendor Model ID
	}

	// Table 4.58: Config Vendor Model Subscription List message parameters
	ConfigVendorModelSubscriptionListMessageParameters struct {
		Status          uint   `bits:"8"`  // Status Code for the requesting message
		ElementAddress  uint   `bits:"16"` // Address of the element
		ModelIdentifier uint   `bits:"32"` // Vendor Model ID
		Addresses       []uint `bits:"8"`  // A block of all addresses from the Subscription List
	}

	// Table 4.59: Config NetKey Add message parameters
	ConfigNetKeyAddMessageParameters struct {
		NetKeyIndex uint   `bits:"12"` // NetKey Index
		Pad         uint   `bits:"4"`  // 0x00
		NetKey      []byte `bits:"8"`  // NetKey
	}

	// Table 4.60: Config NetKey Update message parameters
	ConfigNetKeyUpdateMessageParameters struct {
		NetKeyIndex uint   `bits:"12"` // Index of the NetKey
		Pad         uint   `bits:"4"`  // 0x00
		NetKey      []byte `bits:"8"`  // NetKey
	}

	// Table 4.61: Config NetKey Delete message parameters
	ConfigNetKeyDeleteMessageParameters struct {
		NetKeyIndex uint `bits:"12"` // Index of the NetKey
		Pad         uint `bits:"4"`  // 0x00
	}

	// Table 4.62: Config NetKey Status message parameters
	ConfigNetKeyStatusMessageParameters struct {
		Status      uint `bits:"8"`  // Status Code for the requesting message
		NetKeyIndex uint `bits:"12"` // Index of the NetKey
		Pad         uint `bits:"4"`  // 0x00
	}

	// Table 4.63: Config NetKey List message parameters
	ConfigNetKeyListMessageParameters struct {
		NetKeyIndexes []uint `bits:"12"` // A list of NetKey Indexes known to the node
	}

	// Table 4.64: Config AppKey Add message parameters
	ConfigAppKeyAddMessageParameters struct {
		NetKeyIndex uint   `bits:"12"` // Index of the NetKey
		AppKeyIndex uint   `bits:"12"` //  index of the AppKey
		AppKey      []byte `bits:"8"`  // AppKey value
	}

	// Table 4.65: Config AppKey Update message parameters
	ConfigAppKeyUpdateMessageParameters struct {
		NetKeyIndex uint   `bits:"12"` // Index of the NetKey
		AppKeyIndex uint   `bits:"12"` //  index of the AppKey
		AppKey      []byte `bits:"8"`  // New AppKey value
	}

	// Table 4.66: Config AppKey Delete message parameters
	ConfigAppKeyDeleteMessageParameters struct {
		NetKeyIndex uint `bits:"12"` // Index of the NetKey
		AppKeyIndex uint `bits:"12"` //  index of the AppKey
	}

	// Table 4.67: Config AppKey Status message parameters
	ConfigAppKeyStatusMessageParameters struct {
		Status      uint `bits:"8"`  // Status Code for the requesting message
		NetKeyIndex uint `bits:"12"` // Index of the NetKey
		AppKeyIndex uint `bits:"12"` //  index of the AppKey
	}

	// Table 4.68: Config AppKey Get message parameters
	ConfigAppKeyGetMessageParameters struct {
		NetKeyIndex uint `bits:"12"` // Index of the NetKey
		Pad         uint `bits:"4"`  // 0x00
	}

	// Table 4.69: Config AppKey List message parameters
	ConfigAppKeyListMessageParameters struct {
		Status        uint   `bits:"8"`  // Status Code for the requesting message
		NetKeyIndex   uint   `bits:"12"` // NetKey Index of the NetKey that the AppKeys are bound to
		Pad           uint   `bits:"4"`  // 0x00
		AppKeyIndexes []uint `bits:"12"` // A list of AppKey indexes that are bound to the NetKey identified by NetKeyIndex
	}

	// Table 4.70: Config Node Identity Get message parameters
	ConfigNodeIdentityGetMessageParameters struct {
		NetKeyIndex uint `bits:"12"` // Index of the NetKey
		Pad         uint `bits:"4"`  // 0x00
	}

	// Table 4.71: Config Node Identity Set message parameters
	ConfigNodeIdentitySetMessageParameters struct {
		NetKeyIndex uint `bits:"12"` // Index of the NetKey
		Pad         uint `bits:"4"`  // 0x00
		Identity    uint `bits:"8"`  // New Node Identity state
	}

	// Table 4.72: Config Node Identity Status message parameters
	ConfigNodeIdentityStatusMessageParameters struct {
		Status      uint `bits:"8"`                         // Status Code for the requesting message
		NetKeyIndex uint `bits:"12"`                        // Index of the NetKey
		Pad         uint `bits:"4"`                         // 0x00
		Identity    uint `bits:"8" vt:"NodeIdentityValues"` // Node Identity state
	}

	// Table 4.73: Config Model App Bind message parameters
	ConfigModelAppBindMessageParameters struct {
		ElementAddress  uint `bits:"16"`    // Address of the element
		AppKeyIndex     uint `bits:"12"`    // Index of the AppKey
		Pad             uint `bits:"4"`     // 0x00
		ModelIdentifier uint `bits:"16/32"` // SIG Model ID or Vendor Model ID
	}

	// Table 4.74: Config Model App Unbind message parameters
	ConfigModelAppUnbindMessageParameters struct {
		ElementAddress  uint `bits:"16"`    // Address of the element
		AppKeyIndex     uint `bits:"12"`    // Index of the AppKey
		Pad             uint `bits:"4"`     // 0x00
		ModelIdentifier uint `bits:"16/32"` // SIG Model ID or Vendor Model ID
	}

	// Table 4.75: Config Model App Status message parameters
	ConfigModelAppStatusMessageParameters struct {
		Status          uint `bits:"8"`     // Status Code for the requesting message
		ElementAddress  uint `bits:"16"`    // Address of the element
		AppKeyIndex     uint `bits:"12"`    // Index of the AppKey
		Pad             uint `bits:"4"`     // 0x00
		ModelIdentifier uint `bits:"16/32"` // SIG Model ID or Vendor Model ID
	}

	// Table 4.76: Config SIG Model App Get message parameters
	ConfigSIGModelAppGetMessageParameters struct {
		ElementAddress  uint `bits:"16"` // Address of the element
		ModelIdentifier uint `bits:"16"` // SIG Model ID
	}

	// Table 4.77: Config SIG Model App List message parameters
	ConfigSIGModelAppListMessageParameters struct {
		Status          uint   `bits:"8"`  // Status Code for the requesting message
		ElementAddress  uint   `bits:"16"` // Address of the element
		ModelIdentifier uint   `bits:"16"` // SIG Model ID
		AppKeyIndexes   []uint `bits:"12"` // All AppKey indexes bound to the Model
	}

	// Table 4.78: Config Vendor Model App Get message parameters
	ConfigVendorModelAppGetMessageParameters struct {
		ElementAddress  uint `bits:"16"` // Address of the element
		ModelIdentifier uint `bits:"32"` // Vendor Model ID
	}

	// Table 4.79: Config Vendor Model App List message parameters
	ConfigVendorModelAppListMessageParameters struct {
		Status          uint   `bits:"8"`  // Status Code for the requesting message
		ElementAddress  uint   `bits:"16"` // Address of the element
		ModelIdentifier uint   `bits:"32"` // Vendor Model ID
		AppKeyIndexes   []uint `bits:"12"` // Indexes of all AppKeys bound to the model
	}

	// Table 4.80: Config Friend Set message parameters
	ConfigFriendSetMessageParameters struct {
		Friend uint `bits:"8"` // New Friend state
	}

	// Table 4.81: Config Friend Status message parameters
	ConfigFriendStatusMessageParameters struct {
		Friend uint `bits:"8"` // Friend state
	}

	// Table 4.82: Config Key Refresh Phase Get message parameters
	ConfigKeyRefreshPhaseGetMessageParameters struct {
		NetKeyIndex uint `bits:"12"` // NetKey Index
		Pad         uint `bits:"4"`  // 0x00
	}

	// Table 4.83: Config Key Refresh Phase Set message parameters
	ConfigKeyRefreshPhaseSetMessageParameters struct {
		NetKeyIndex uint `bits:"12"` // NetKey Index
		Pad         uint `bits:"4"`  // 0x00
		Transition  uint `bits:"8"`  // New Key Refresh Phase Transition
	}

	// Table 4.84: Config Key Refresh Phase Status message parameters
	ConfigKeyRefreshPhaseStatusMessageParameters struct {
		Status      uint `bits:"8"`  // Status Code for the requesting message
		NetKeyIndex uint `bits:"12"` // NetKey Index
		Pad         uint `bits:"4"`  // 0x00
		Phase       uint `bits:"8"`  // Key Refresh Phase State
	}

	// Table 4.85: Config Heartbeat Publication Set message parameters
	ConfigHeartbeatPublicationSetMessageParameters struct {
		Destination uint `bits:"16"` // Destination address for Heartbeat messages
		CountLog    uint `bits:"8"`  // Number of Heartbeat messages to be sent
		PeriodLog   uint `bits:"8"`  // Period for sending Heartbeat messages
		TTL         uint `bits:"8"`  // TTL to be used when sending Heartbeat messages
		Features    uint `bits:"16"` // Bit field indicating features that trigger Heartbeat messages when changed
		NetKeyIndex uint `bits:"12"` // NetKey Index
		Pad         uint `bits:"4"`  // 0x00
	}

	// Table 4.86: Config Heartbeat Publication Status message parameters
	ConfigHeartbeatPublicationStatusMessageParameters struct {
		Status      uint `bits:"8"`  // Status Code for the requesting message
		Destination uint `bits:"16"` // Destination address for Heartbeat messages
		CountLog    uint `bits:"8"`  // Number of Heartbeat messages remaining to be sent
		PeriodLog   uint `bits:"8"`  // Period for sending Heartbeat messages
		TTL         uint `bits:"8"`  // TTL to be used when sending Heartbeat messages
		Features    uint `bits:"16"` // Bit field indicating features that trigger Heartbeat messages when changed
		NetKeyIndex uint `bits:"12"` // NetKey Index
		Pad         uint `bits:"4"`  // 0x00
	}

	// Table 4.87: Config Heartbeat Subscription Set message parameters
	ConfigHeartbeatSubscriptionSetMessageParameters struct {
		Source      uint `bits:"16"` // Source address for Heartbeat messages
		Destination uint `bits:"16"` // Destination address for Heartbeat messages
		PeriodLog   uint `bits:"8"`  // Period for receiving Heartbeat messages
	}

	// Table 4.88: Config Heartbeat Subscription Status message parameters
	ConfigHeartbeatSubscriptionStatusMessageParameters struct {
		Status      uint `bits:"8"`  // Status Code for the requesting message
		Source      uint `bits:"16"` // Source address for Heartbeat messages
		Destination uint `bits:"16"` // Destination address for Heartbeat messages
		PeriodLog   uint `bits:"8"`  // Remaining Period for processing Heartbeat messages
		CountLog    uint `bits:"8"`  // Number of Heartbeat messages received
		MinHops     uint `bits:"8"`  // Minimum hops when receiving Heartbeat messages
		MaxHops     uint `bits:"8"`  // Maximum hops when receiving Heartbeat messages
	}

	// Table 4.89: Config Low Power Node PollTimeout Get message parameters
	ConfigLowPowerNodePollTimeoutGetMessageParameters struct {
		LPNAddress uint `bits:"16"` // The unicast address of the Low Power node
	}

	// Table 4.90: Config Low Power Node PollTimeout Status message parameters
	ConfigLowPowerNodePollTimeoutStatusMessageParameters struct {
		LPNAddress  uint `bits:"16"` // The unicast address of the Low Power node
		PollTimeout uint `bits:"24"` // The current value of the PollTimeout timer of the Low Power node
	}

	// Table 4.91: Config Network Transmit Set message parameters
	ConfigNetworkTransmitSetMessageParameters struct {
		NetworkTransmitCount         uint `bits:"3"` // Number of transmissions for each Network PDU originating from the node
		NetworkTransmitIntervalSteps uint `bits:"5"` // Number of 10-millisecond steps between transmissions
	}

	// Table 4.92: Config Network Transmit Status message parameters
	ConfigNetworkTransmitStatusMessageParameters struct {
		NetworkTransmitCount         uint `bits:"3"` // Number of transmissions for each Network PDU originating from the node
		NetworkTransmitIntervalSteps uint `bits:"5"` // Number of 10-millisecond steps between transmissions
	}

	// Table 4.93: Health Current Status message parameters
	HealthCurrentStatusMessageParameters struct {
		TestID     uint   `bits:"8"`  // Identifier of a most recently performed test
		CompanyID  uint   `bits:"16"` // 16-bit Bluetooth assigned Company Identifier
		FaultArray []uint `bits:"8"`  // The FaultArray field contains a sequence of 1-octet fault values
	}

	// Table 4.94: Health Fault Get message parameters
	HealthFaultGetMessageParameters struct {
		CompanyID uint `bits:"16"` // 16-bit Bluetooth assigned Company Identifier
	}

	// Table 4.95: Health Fault Clear Unacknowledged message parameters
	HealthFaultClearUnacknowledgedMessageParameters struct {
		CompanyID uint `bits:"16"` // 16-bit Bluetooth assigned Company Identifier
	}

	// Table 4.96: Health Fault Clear message parameters
	HealthFaultClearMessageParameters struct {
		CompanyID uint `bits:"16"` // 16-bit Bluetooth assigned Company Identifier
	}

	// Table 4.97: Health Fault Test message parameters
	HealthFaultTestMessageParameters struct {
		TestID    uint `bits:"8"`  // Identifier of a specific test to be performed
		CompanyID uint `bits:"16"` // 16-bit Bluetooth assigned Company Identifier
	}

	// Table 4.98: Health Fault Test Unacknowledged message parameters
	HealthFaultTestUnacknowledgedMessageParameters struct {
		TestID    uint `bits:"8"`  // Identifier of a specific test to be performed
		CompanyID uint `bits:"16"` // 16-bit Bluetooth assigned Company Identifier
	}

	// Table 4.99: Health Fault Status message parameters
	HealthFaultStatusMessageParameters struct {
		TestID     uint `bits:"8"`  // Identifier of a most recently performed test
		CompanyID  uint `bits:"16"` // 16-bit Bluetooth assigned Company Identifier
		FaultArray uint `bits:"N"`  // The FaultArray field contains a sequence of 1-octet fault values
	}

	// Table 4.100: Health Period Set Unacknowledged message parameters
	HealthPeriodSetUnacknowledgedMessageParameters struct {
		FastPeriodDivisor uint `bits:"8"` // Divider for the Publish Period. Modified Publish Period is used for sending Current Health Status messages when there are active faults to communicate
	}

	// Table 4.101: Health Period Set message parameters
	HealthPeriodSetMessageParameters struct {
		FastPeriodDivisor uint `bits:"8"` // Divider for the Publish Period. Modified Publish Period is used for sending Current Health Status messages when there are active faults to communicate
	}

	// Table 4.102: Health Period Status message parameters
	HealthPeriodStatusMessageParameters struct {
		FastPeriodDivisor uint `bits:"8"` // Divider for the Publish Period. Modified Publish Period is used for sending Current Health Status messages when there are active faults to communicate
	}

	// Table 4.103: Health Attention Set message parameters
	HealthAttentionSetMessageParameters struct {
		Attention uint `bits:"8"` // Value of the Attention Timer state
	}

	// Table 4.104: Health Attention Set Unacknowledged message parameters
	HealthAttentionSetUnacknowledgedMessageParameters struct {
		Attention uint `bits:"8"` // Value of the Attention Timer state
	}

	// Table 4.105: Attention Status message parameters
	AttentionStatusMessageParameters struct {
		Attention uint `bits:"8"` // Value of the Attention Timer state
	}

	// Table 5.15: Provisioning Invite PDU parameters format
	ProvisioningInvitePDUParametersFormat struct {
		AttentionDuration uint `bits:"8"` // Attention Timer state (See Section 4.2.9)
	}

	// Table 5.16: Provisioning capabilities PDU parameters format
	ProvisioningCapabilitiesPDUParametersFormat struct {
		NumberOfElements uint `bits:"8"`  // Number of elements supported by the device (Table 5.17)
		Algorithms       uint `bits:"16"` // Supported algorithms and other capabilities (see Table 5.18)
		PublicKeyType    uint `bits:"8"`  // Supported public key types (see Table 5.19)
		StaticOOBType    uint `bits:"8"`  // Supported static OOB Types (see Table 5.20)
		OutputOOBSize    uint `bits:"8"`  // Maximum size of Output OOB supported (see Table 5.21)
		OutputOOBAction  uint `bits:"16"` // Supported Output OOB Actions (see Table 5.22)
		InputOOBSize     uint `bits:"8"`  // Maximum size in octets of Input OOB supported (see Table 5.23)
		InputOOBAction   uint `bits:"16"` // Supported Input OOB Actions (see Table 5.24)
	}

	// Table 5.25: Provisioning Start PDU parameters format
	ProvisioningStartPDUParametersFormat struct {
		Algorithm            uint `bits:"8"` // The algorithm used for provisioning (see Table 5.26)
		PublicKey            uint `bits:"8"` // Public Key used (see Table 5.27)
		AuthenticationMethod uint `bits:"8"` // Authentication Method used (see Table 5.28)
		AuthenticationAction uint `bits:"8"` // Selected Output OOB Action (see Table 5.29) or Input OOB Action (see Table 5.31) or 0x00
		AuthenticationSize   uint `bits:"8"` // Size of the Output OOB used (see Table 5.30) or size of the Input OOB used (see Table 5.32) or 0x00
	}

	// Table 5.33: Provisioning Public Key PDU Parameters Format
	ProvisioningPublicKeyPDUParametersFormat struct {
		PublicKeyX uint `bits:"256"` // The X component of public key for the FIPS P-256 algorithm
		PublicKeyY uint `bits:"256"` // The Y component of public key for the FIPS P-256 algorithm
	}

	// Table 5.34: Provisioning Confirmation PDU Parameters Format
	ProvisioningConfirmationPDUParametersFormat struct {
		Confirmation uint `bits:"128"` // The values exchanged so far including the OOB Authentication value
	}

	// Table 5.35: Provisioning Random PDU parameters format
	ProvisioningRandomPDUParametersFormat struct {
		Random uint `bits:"128"` // The final input to the confirmation
	}

	// Table 5.36: Provisioning Data PDU parameters format
	ProvisioningDataPDUParametersFormat struct {
		EncryptedProvisioningData uint `bits:"200"` // An encrypted and authenticated network key, NetKey Index, Key Refresh Flag, IV Update Flag, current value of the IV Index, and unicast address of the primary element (see Section 5.4.2.5)
		ProvisioningDataMIC       uint `bits:"64"`  // PDU Integrity Check value
	}

	// Table 5.37: Provisioning Failed PDU parameters format
	ProvisioningFailedPDUParametersFormat struct {
		ErrorCode uint `bits:"8"` // This represents a specific error in the provisioning protocol encountered by a device
	}

	// Table 5.39: Provisioning data format
	ProvisioningDataFormat struct {
		NetworkKey     uint `bits:"128"` // NetKey
		KeyIndex       uint `bits:"12"`  // Index of the NetKey
		Pad            uint `bits:"4"`   // 0x00
		Flags          uint `bits:"8"`   // Flags bitmask
		IVIndex        uint `bits:"32"`  // Current value of the IV Index
		UnicastAddress uint `bits:"16"`  // Unicast address of the primary element
	}

	// Table 6.1: Proxy PDU format
	ProxyPDUFormat struct {
		SAR         uint   `bits:"2"` // Message segmentation and reassembly information
		MessageType uint   `bits:"6"` // Type of message contained in the PDU
		Data        []uint `bits:"8"` // Full message or message segment
	}

	// Table 6.4: Proxy TransportPDU field format
	ProxyTransportPDUFieldFormat struct {
		Opcode uint `bits:"8"` // Message opcode
	}

	// Table 6.6: Set Filter Type Message Format
	SetFilterTypeMessageFormat struct {
		FilterType uint `bits:"8"` // The proxy filter type.
	}

	// Table 6.8: Add Addresses to Filter message format
	AddAddressesToFilterMessageFormat struct {
		AddressArray []uint `bits:"16"` // List of addresses where N is the number of addresses in this message.
	}

	// Table 6.9: Remove Addresses from Filter message format
	RemoveAddressesFromFilterMessageFormat struct {
		AddressArray []uint `bits:"16"` // List of addresses where N is the number of addresses in this message.
	}

	// Table 6.10: Filter Status message format
	FilterStatusMessageFormat struct {
		FilterType uint `bits:"8"`  // White list or black list.
		ListSize   uint `bits:"16"` // Number of addresses in the proxy filter list.
	}

	// Table 7.3: Service Data for Mesh Provisioning Service
	ServiceDataForMeshProvisioningService struct {
		DeviceUUID     uint `bits:"128"` // See Section 3.10.3
		OOBInformation uint `bits:"16"`  // See Table 3.54
	}

	// Table 7.7: Service Data for Mesh Proxy Service
	ServiceDataForMeshProxyService struct {
		IdentificationType uint `bits:"8"` // See Table 7.8
	}

	// Table 7.9: Service Data for Mesh Proxy Service with Network ID
	ServiceDataForMeshProxyServiceWithNetworkID struct {
		IdentificationType uint `bits:"8"`  // 0x00 (Network ID type)
		NetworkID          uint `bits:"64"` // Identifies the network
	}

	// Table 7.10: Service Data for Mesh Proxy Service with Node Identity
	ServiceDataForMeshProxyServiceWithNodeIdentity struct {
		IdentificationType uint `bits:"8"`  // 0x01 (Node Identity type)
		Hash               uint `bits:"64"` // Function of the included random number and identity information.
		Random             uint `bits:"64"` // 64-bit random number
	}
)

var (
	vv = map[Range]string{
		Range{0, 1}: "",
	}
)
