package gtp

var (
    GTPV2_IE_IMSI = NewGTPv2IE("GTPV2_IE_IMSI", 1, 0, "International Mobile Subscriber Identity (IMSI)")
    GTPV2_IE_CAUSE = NewGTPv2IE("GTPV2_IE_CAUSE", 2, 0, "Cause")
    GTPV2_IE_RECOVERY = NewGTPv2IE("GTPV2_IE_RECOVERY", 3, 0, "Recovery (Restart Counter)")
    GTPV2_IE_STN_SR = NewGTPv2IE("GTPV2_IE_STN_SR", 51, 0, "STN-SR")
    GTPV2_IE_APN = NewGTPv2IE("GTPV2_IE_APN", 71, 0, "Access Point Name (APN)")
    GTPV2_IE_AMBR = NewGTPv2IE("GTPV2_IE_AMBR", 72, 8, "Aggregate Maximum Bit Rate (AMBR)")
    GTPV2_IE_EBI = NewGTPv2IE("GTPV2_IE_EBI", 73, 1, "EPS Bearer ID (EBI)")
    GTPV2_IE_IP_ADDRESS = NewGTPv2IE("GTPV2_IE_IP_ADDRESS", 74, 0, "IP Address")
    GTPV2_IE_MEI = NewGTPv2IE("GTPV2_IE_MEI", 75, 0, "Mobile Equipment Identity")
    GTPV2_IE_MSISDN = NewGTPv2IE("GTPV2_IE_MSISDN", 76, 0, "MSISDN")
)
