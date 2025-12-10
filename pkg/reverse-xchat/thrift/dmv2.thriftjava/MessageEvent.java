package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import android.gov.nist.core.net.C0009a;
import android.gov.nist.javax.sip.clientauthutils.C0026b;
import android.gov.nist.javax.sip.header.C0031b;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.ThriftException;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import com.google.ads.interactivemedia.p012v3.impl.data.C12781a;
import java.io.IOException;
import kotlin.Metadata;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;
import tv.periscope.android.api.Constants;

@Metadata(m64929d1 = {"\u0000J\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\u000e\n\u0002\b\u0006\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0010\u000b\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\u0015\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0002\b\u000b\b\u0086\b\u0018\u0000 82\u00020\u0001:\u000298Bu\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0004\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0005\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0006\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0007\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\b\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\n\u001a\u0004\u0018\u00010\t\u0012\b\u0010\f\u001a\u0004\u0018\u00010\u000b\u0012\b\u0010\u000e\u001a\u0004\u0018\u00010\r\u0012\b\u0010\u000f\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0011\u001a\u0004\u0018\u00010\u0010¢\u0006\u0004\b\u0012\u0010\u0013J\u0017\u0010\u0017\u001a\u00020\u00162\u0006\u0010\u0015\u001a\u00020\u0014H\u0016¢\u0006\u0004\b\u0017\u0010\u0018J\u0012\u0010\u0019\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0019\u0010\u001aJ\u0012\u0010\u001b\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u001b\u0010\u001aJ\u0012\u0010\u001c\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u001c\u0010\u001aJ\u0012\u0010\u001d\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u001d\u0010\u001aJ\u0012\u0010\u001e\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u001e\u0010\u001aJ\u0012\u0010\u001f\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u001f\u0010\u001aJ\u0012\u0010 \u001a\u0004\u0018\u00010\tHÆ\u0003¢\u0006\u0004\b \u0010!J\u0012\u0010\"\u001a\u0004\u0018\u00010\u000bHÆ\u0003¢\u0006\u0004\b\"\u0010#J\u0012\u0010$\u001a\u0004\u0018\u00010\rHÆ\u0003¢\u0006\u0004\b$\u0010%J\u0012\u0010&\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b&\u0010\u001aJ\u0012\u0010'\u001a\u0004\u0018\u00010\u0010HÆ\u0003¢\u0006\u0004\b'\u0010(J\u0094\u0001\u0010)\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0004\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0005\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0006\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0007\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\b\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\n\u001a\u0004\u0018\u00010\t2\n\b\u0002\u0010\f\u001a\u0004\u0018\u00010\u000b2\n\b\u0002\u0010\u000e\u001a\u0004\u0018\u00010\r2\n\b\u0002\u0010\u000f\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0011\u001a\u0004\u0018\u00010\u0010HÆ\u0001¢\u0006\u0004\b)\u0010*J\u0010\u0010+\u001a\u00020\u0002HÖ\u0001¢\u0006\u0004\b+\u0010\u001aJ\u0010\u0010-\u001a\u00020,HÖ\u0001¢\u0006\u0004\b-\u0010.J\u001a\u00101\u001a\u00020\u00102\b\u00100\u001a\u0004\u0018\u00010/HÖ\u0003¢\u0006\u0004\b1\u00102R\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u00103R\u0016\u0010\u0004\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0004\u00103R\u0016\u0010\u0005\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0005\u00103R\u0016\u0010\u0006\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0006\u00103R\u0016\u0010\u0007\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0007\u00103R\u0016\u0010\b\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\b\u00103R\u0016\u0010\n\u001a\u0004\u0018\u00010\t8\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\n\u00104R\u0016\u0010\f\u001a\u0004\u0018\u00010\u000b8\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\f\u00105R\u0016\u0010\u000e\u001a\u0004\u0018\u00010\r8\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u000e\u00106R\u0016\u0010\u000f\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u000f\u00103R\u0016\u0010\u0011\u001a\u0004\u0018\u00010\u00108\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0011\u00107¨\u0006:"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEvent;", "Lcom/bendb/thrifty/a;", "", "sequence_id", "message_id", "sender_id", "conversation_id", "conversation_token", "created_at_msec", "Lcom/x/dmv2/thriftjava/MessageEventDetail;", "detail", "Lcom/x/dmv2/thriftjava/MessageEventRelaySource;", "relay_source", "Lcom/x/dmv2/thriftjava/MessageEventSignature;", "message_event_signature", "previous_sequence_id", "", "is_trusted", "<init>", "(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Lcom/x/dmv2/thriftjava/MessageEventDetail;Lcom/x/dmv2/thriftjava/MessageEventRelaySource;Lcom/x/dmv2/thriftjava/MessageEventSignature;Ljava/lang/String;Ljava/lang/Boolean;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/String;", "component2", "component3", "component4", "component5", "component6", "component7", "()Lcom/x/dmv2/thriftjava/MessageEventDetail;", "component8", "()Lcom/x/dmv2/thriftjava/MessageEventRelaySource;", "component9", "()Lcom/x/dmv2/thriftjava/MessageEventSignature;", "component10", "component11", "()Ljava/lang/Boolean;", "copy", "(Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Lcom/x/dmv2/thriftjava/MessageEventDetail;Lcom/x/dmv2/thriftjava/MessageEventRelaySource;Lcom/x/dmv2/thriftjava/MessageEventSignature;Ljava/lang/String;Ljava/lang/Boolean;)Lcom/x/dmv2/thriftjava/MessageEvent;", "toString", "", "hashCode", "()I", "", "other", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/String;", "Lcom/x/dmv2/thriftjava/MessageEventDetail;", "Lcom/x/dmv2/thriftjava/MessageEventRelaySource;", "Lcom/x/dmv2/thriftjava/MessageEventSignature;", "Ljava/lang/Boolean;", "Companion", "MessageEventAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class MessageEvent implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final String conversation_id;

    @JvmField
    @InterfaceC88465b
    public final String conversation_token;

    @JvmField
    @InterfaceC88465b
    public final String created_at_msec;

    @JvmField
    @InterfaceC88465b
    public final MessageEventDetail detail;

    @JvmField
    @InterfaceC88465b
    public final Boolean is_trusted;

    @JvmField
    @InterfaceC88465b
    public final MessageEventSignature message_event_signature;

    @JvmField
    @InterfaceC88465b
    public final String message_id;

    @JvmField
    @InterfaceC88465b
    public final String previous_sequence_id;

    @JvmField
    @InterfaceC88465b
    public final MessageEventRelaySource relay_source;

    @JvmField
    @InterfaceC88465b
    public final String sender_id;

    @JvmField
    @InterfaceC88465b
    public final String sequence_id;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new MessageEventAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEvent$MessageEventAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/MessageEvent;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/MessageEvent;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/MessageEvent;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class MessageEventAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public MessageEvent m85658read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            String string = null;
            String string2 = null;
            String string3 = null;
            String string4 = null;
            String string5 = null;
            String string6 = null;
            MessageEventDetail messageEventDetail = null;
            MessageEventRelaySource messageEventRelaySource = null;
            MessageEventSignature messageEventSignature = null;
            String string7 = null;
            Boolean boolValueOf = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b != 0) {
                    switch (c11265cMo14127V2.f38393b) {
                        case 1:
                            if (b != 11) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                string = protocol.readString();
                                break;
                            }
                        case 2:
                            if (b != 11) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                string2 = protocol.readString();
                                break;
                            }
                        case 3:
                            if (b != 11) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                string3 = protocol.readString();
                                break;
                            }
                        case 4:
                            if (b != 11) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                string4 = protocol.readString();
                                break;
                            }
                        case 5:
                            if (b != 11) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                string5 = protocol.readString();
                                break;
                            }
                        case 6:
                            if (b != 11) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                string6 = protocol.readString();
                                break;
                            }
                        case 7:
                            if (b != 12) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                messageEventDetail = (MessageEventDetail) MessageEventDetail.ADAPTER.read(protocol);
                                break;
                            }
                        case 8:
                            if (b != 8) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                int iMo14132c4 = protocol.mo14132c4();
                                MessageEventRelaySource messageEventRelaySourceFindByValue = MessageEventRelaySource.INSTANCE.findByValue(iMo14132c4);
                                if (messageEventRelaySourceFindByValue == null) {
                                    throw new ThriftException(ThriftException.EnumC11260b.PROTOCOL_ERROR, C0031b.m45c(iMo14132c4, "Unexpected value for enum type MessageEventRelaySource: "));
                                }
                                messageEventRelaySource = messageEventRelaySourceFindByValue;
                                break;
                            }
                        case 9:
                            if (b != 12) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                messageEventSignature = (MessageEventSignature) MessageEventSignature.ADAPTER.read(protocol);
                                break;
                            }
                        case 10:
                            if (b != 11) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                string7 = protocol.readString();
                                break;
                            }
                        case 11:
                            if (b != 2) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                boolValueOf = Boolean.valueOf(protocol.readBool());
                                break;
                            }
                        default:
                            C11272a.m14141a(protocol, b);
                            break;
                    }
                } else {
                    return new MessageEvent(string, string2, string3, string4, string5, string6, messageEventDetail, messageEventRelaySource, messageEventSignature, string7, boolValueOf);
                }
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a MessageEvent struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("MessageEvent");
            if (struct.sequence_id != null) {
                protocol.mo14136v3("sequence_id", 1, (byte) 11);
                protocol.mo14137w0(struct.sequence_id);
            }
            if (struct.message_id != null) {
                protocol.mo14136v3("message_id", 2, (byte) 11);
                protocol.mo14137w0(struct.message_id);
            }
            if (struct.sender_id != null) {
                protocol.mo14136v3("sender_id", 3, (byte) 11);
                protocol.mo14137w0(struct.sender_id);
            }
            if (struct.conversation_id != null) {
                protocol.mo14136v3("conversation_id", 4, (byte) 11);
                protocol.mo14137w0(struct.conversation_id);
            }
            if (struct.conversation_token != null) {
                protocol.mo14136v3("conversation_token", 5, (byte) 11);
                protocol.mo14137w0(struct.conversation_token);
            }
            if (struct.created_at_msec != null) {
                protocol.mo14136v3("created_at_msec", 6, (byte) 11);
                protocol.mo14137w0(struct.created_at_msec);
            }
            if (struct.detail != null) {
                protocol.mo14136v3("detail", 7, (byte) 12);
                MessageEventDetail.ADAPTER.write(protocol, struct.detail);
            }
            if (struct.relay_source != null) {
                protocol.mo14136v3("relay_source", 8, (byte) 8);
                protocol.mo14122C2(struct.relay_source.value);
            }
            if (struct.message_event_signature != null) {
                protocol.mo14136v3("message_event_signature", 9, (byte) 12);
                MessageEventSignature.ADAPTER.write(protocol, struct.message_event_signature);
            }
            if (struct.previous_sequence_id != null) {
                protocol.mo14136v3("previous_sequence_id", 10, (byte) 11);
                protocol.mo14137w0(struct.previous_sequence_id);
            }
            if (struct.is_trusted != null) {
                protocol.mo14136v3("is_trusted", 11, (byte) 2);
                protocol.mo14125P1(struct.is_trusted.booleanValue());
            }
            protocol.mo14134i0();
        }
    }

    public MessageEvent(@InterfaceC88465b String str, @InterfaceC88465b String str2, @InterfaceC88465b String str3, @InterfaceC88465b String str4, @InterfaceC88465b String str5, @InterfaceC88465b String str6, @InterfaceC88465b MessageEventDetail messageEventDetail, @InterfaceC88465b MessageEventRelaySource messageEventRelaySource, @InterfaceC88465b MessageEventSignature messageEventSignature, @InterfaceC88465b String str7, @InterfaceC88465b Boolean bool) {
        this.sequence_id = str;
        this.message_id = str2;
        this.sender_id = str3;
        this.conversation_id = str4;
        this.conversation_token = str5;
        this.created_at_msec = str6;
        this.detail = messageEventDetail;
        this.relay_source = messageEventRelaySource;
        this.message_event_signature = messageEventSignature;
        this.previous_sequence_id = str7;
        this.is_trusted = bool;
    }

    public static /* synthetic */ MessageEvent copy$default(MessageEvent messageEvent, String str, String str2, String str3, String str4, String str5, String str6, MessageEventDetail messageEventDetail, MessageEventRelaySource messageEventRelaySource, MessageEventSignature messageEventSignature, String str7, Boolean bool, int i, Object obj) {
        return messageEvent.copy((i & 1) != 0 ? messageEvent.sequence_id : str, (i & 2) != 0 ? messageEvent.message_id : str2, (i & 4) != 0 ? messageEvent.sender_id : str3, (i & 8) != 0 ? messageEvent.conversation_id : str4, (i & 16) != 0 ? messageEvent.conversation_token : str5, (i & 32) != 0 ? messageEvent.created_at_msec : str6, (i & 64) != 0 ? messageEvent.detail : messageEventDetail, (i & 128) != 0 ? messageEvent.relay_source : messageEventRelaySource, (i & 256) != 0 ? messageEvent.message_event_signature : messageEventSignature, (i & 512) != 0 ? messageEvent.previous_sequence_id : str7, (i & Constants.BITS_PER_KILOBIT) != 0 ? messageEvent.is_trusted : bool);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final String getSequence_id() {
        return this.sequence_id;
    }

    @InterfaceC88465b
    /* renamed from: component10, reason: from getter */
    public final String getPrevious_sequence_id() {
        return this.previous_sequence_id;
    }

    @InterfaceC88465b
    /* renamed from: component11, reason: from getter */
    public final Boolean getIs_trusted() {
        return this.is_trusted;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final String getMessage_id() {
        return this.message_id;
    }

    @InterfaceC88465b
    /* renamed from: component3, reason: from getter */
    public final String getSender_id() {
        return this.sender_id;
    }

    @InterfaceC88465b
    /* renamed from: component4, reason: from getter */
    public final String getConversation_id() {
        return this.conversation_id;
    }

    @InterfaceC88465b
    /* renamed from: component5, reason: from getter */
    public final String getConversation_token() {
        return this.conversation_token;
    }

    @InterfaceC88465b
    /* renamed from: component6, reason: from getter */
    public final String getCreated_at_msec() {
        return this.created_at_msec;
    }

    @InterfaceC88465b
    /* renamed from: component7, reason: from getter */
    public final MessageEventDetail getDetail() {
        return this.detail;
    }

    @InterfaceC88465b
    /* renamed from: component8, reason: from getter */
    public final MessageEventRelaySource getRelay_source() {
        return this.relay_source;
    }

    @InterfaceC88465b
    /* renamed from: component9, reason: from getter */
    public final MessageEventSignature getMessage_event_signature() {
        return this.message_event_signature;
    }

    @InterfaceC88464a
    public final MessageEvent copy(@InterfaceC88465b String sequence_id, @InterfaceC88465b String message_id, @InterfaceC88465b String sender_id, @InterfaceC88465b String conversation_id, @InterfaceC88465b String conversation_token, @InterfaceC88465b String created_at_msec, @InterfaceC88465b MessageEventDetail detail, @InterfaceC88465b MessageEventRelaySource relay_source, @InterfaceC88465b MessageEventSignature message_event_signature, @InterfaceC88465b String previous_sequence_id, @InterfaceC88465b Boolean is_trusted) {
        return new MessageEvent(sequence_id, message_id, sender_id, conversation_id, conversation_token, created_at_msec, detail, relay_source, message_event_signature, previous_sequence_id, is_trusted);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof MessageEvent)) {
            return false;
        }
        MessageEvent messageEvent = (MessageEvent) other;
        return Intrinsics.m65267c(this.sequence_id, messageEvent.sequence_id) && Intrinsics.m65267c(this.message_id, messageEvent.message_id) && Intrinsics.m65267c(this.sender_id, messageEvent.sender_id) && Intrinsics.m65267c(this.conversation_id, messageEvent.conversation_id) && Intrinsics.m65267c(this.conversation_token, messageEvent.conversation_token) && Intrinsics.m65267c(this.created_at_msec, messageEvent.created_at_msec) && Intrinsics.m65267c(this.detail, messageEvent.detail) && this.relay_source == messageEvent.relay_source && Intrinsics.m65267c(this.message_event_signature, messageEvent.message_event_signature) && Intrinsics.m65267c(this.previous_sequence_id, messageEvent.previous_sequence_id) && Intrinsics.m65267c(this.is_trusted, messageEvent.is_trusted);
    }

    public int hashCode() {
        String str = this.sequence_id;
        int iHashCode = (str == null ? 0 : str.hashCode()) * 31;
        String str2 = this.message_id;
        int iHashCode2 = (iHashCode + (str2 == null ? 0 : str2.hashCode())) * 31;
        String str3 = this.sender_id;
        int iHashCode3 = (iHashCode2 + (str3 == null ? 0 : str3.hashCode())) * 31;
        String str4 = this.conversation_id;
        int iHashCode4 = (iHashCode3 + (str4 == null ? 0 : str4.hashCode())) * 31;
        String str5 = this.conversation_token;
        int iHashCode5 = (iHashCode4 + (str5 == null ? 0 : str5.hashCode())) * 31;
        String str6 = this.created_at_msec;
        int iHashCode6 = (iHashCode5 + (str6 == null ? 0 : str6.hashCode())) * 31;
        MessageEventDetail messageEventDetail = this.detail;
        int iHashCode7 = (iHashCode6 + (messageEventDetail == null ? 0 : messageEventDetail.hashCode())) * 31;
        MessageEventRelaySource messageEventRelaySource = this.relay_source;
        int iHashCode8 = (iHashCode7 + (messageEventRelaySource == null ? 0 : messageEventRelaySource.hashCode())) * 31;
        MessageEventSignature messageEventSignature = this.message_event_signature;
        int iHashCode9 = (iHashCode8 + (messageEventSignature == null ? 0 : messageEventSignature.hashCode())) * 31;
        String str7 = this.previous_sequence_id;
        int iHashCode10 = (iHashCode9 + (str7 == null ? 0 : str7.hashCode())) * 31;
        Boolean bool = this.is_trusted;
        return iHashCode10 + (bool != null ? bool.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        String str = this.sequence_id;
        String str2 = this.message_id;
        String str3 = this.sender_id;
        String str4 = this.conversation_id;
        String str5 = this.conversation_token;
        String str6 = this.created_at_msec;
        MessageEventDetail messageEventDetail = this.detail;
        MessageEventRelaySource messageEventRelaySource = this.relay_source;
        MessageEventSignature messageEventSignature = this.message_event_signature;
        String str7 = this.previous_sequence_id;
        Boolean bool = this.is_trusted;
        StringBuilder sbM11b = C0009a.m11b("MessageEvent(sequence_id=", str, ", message_id=", str2, ", sender_id=");
        C0026b.m37b(sbM11b, str3, ", conversation_id=", str4, ", conversation_token=");
        C0026b.m37b(sbM11b, str5, ", created_at_msec=", str6, ", detail=");
        sbM11b.append(messageEventDetail);
        sbM11b.append(", relay_source=");
        sbM11b.append(messageEventRelaySource);
        sbM11b.append(", message_event_signature=");
        sbM11b.append(messageEventSignature);
        sbM11b.append(", previous_sequence_id=");
        sbM11b.append(str7);
        sbM11b.append(", is_trusted=");
        return C12781a.m16257a(sbM11b, bool, Separators.RPAREN);
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}