package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import android.gov.nist.javax.sip.header.C0031b;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.ThriftException;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import com.x.composer.sensitivemedia.C68546u;
import java.io.IOException;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;
import kotlin.Metadata;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u0000b\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\u000e\n\u0000\n\u0002\u0010 \n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\u0013\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\u000b\b\u0086\b\u0018\u0000 92\u00020\u0001:\u0002:9Bi\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\u000e\u0010\u0006\u001a\n\u0012\u0004\u0012\u00020\u0005\u0018\u00010\u0004\u0012\u000e\u0010\b\u001a\n\u0012\u0004\u0012\u00020\u0007\u0018\u00010\u0004\u0012\b\u0010\n\u001a\u0004\u0018\u00010\t\u0012\b\u0010\f\u001a\u0004\u0018\u00010\u000b\u0012\b\u0010\u000e\u001a\u0004\u0018\u00010\r\u0012\b\u0010\u0010\u001a\u0004\u0018\u00010\u000f\u0012\u000e\u0010\u0012\u001a\n\u0012\u0004\u0012\u00020\u0011\u0018\u00010\u0004¢\u0006\u0004\b\u0013\u0010\u0014J\u0017\u0010\u0018\u001a\u00020\u00172\u0006\u0010\u0016\u001a\u00020\u0015H\u0016¢\u0006\u0004\b\u0018\u0010\u0019J\u0012\u0010\u001a\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u001a\u0010\u001bJ\u0018\u0010\u001c\u001a\n\u0012\u0004\u0012\u00020\u0005\u0018\u00010\u0004HÆ\u0003¢\u0006\u0004\b\u001c\u0010\u001dJ\u0018\u0010\u001e\u001a\n\u0012\u0004\u0012\u00020\u0007\u0018\u00010\u0004HÆ\u0003¢\u0006\u0004\b\u001e\u0010\u001dJ\u0012\u0010\u001f\u001a\u0004\u0018\u00010\tHÆ\u0003¢\u0006\u0004\b\u001f\u0010 J\u0012\u0010!\u001a\u0004\u0018\u00010\u000bHÆ\u0003¢\u0006\u0004\b!\u0010\"J\u0012\u0010#\u001a\u0004\u0018\u00010\rHÆ\u0003¢\u0006\u0004\b#\u0010$J\u0012\u0010%\u001a\u0004\u0018\u00010\u000fHÆ\u0003¢\u0006\u0004\b%\u0010&J\u0018\u0010'\u001a\n\u0012\u0004\u0012\u00020\u0011\u0018\u00010\u0004HÆ\u0003¢\u0006\u0004\b'\u0010\u001dJ\u0082\u0001\u0010(\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\u0010\b\u0002\u0010\u0006\u001a\n\u0012\u0004\u0012\u00020\u0005\u0018\u00010\u00042\u0010\b\u0002\u0010\b\u001a\n\u0012\u0004\u0012\u00020\u0007\u0018\u00010\u00042\n\b\u0002\u0010\n\u001a\u0004\u0018\u00010\t2\n\b\u0002\u0010\f\u001a\u0004\u0018\u00010\u000b2\n\b\u0002\u0010\u000e\u001a\u0004\u0018\u00010\r2\n\b\u0002\u0010\u0010\u001a\u0004\u0018\u00010\u000f2\u0010\b\u0002\u0010\u0012\u001a\n\u0012\u0004\u0012\u00020\u0011\u0018\u00010\u0004HÆ\u0001¢\u0006\u0004\b(\u0010)J\u0010\u0010*\u001a\u00020\u0002HÖ\u0001¢\u0006\u0004\b*\u0010\u001bJ\u0010\u0010,\u001a\u00020+HÖ\u0001¢\u0006\u0004\b,\u0010-J\u001a\u00101\u001a\u0002002\b\u0010/\u001a\u0004\u0018\u00010.HÖ\u0003¢\u0006\u0004\b1\u00102R\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u00103R\u001c\u0010\u0006\u001a\n\u0012\u0004\u0012\u00020\u0005\u0018\u00010\u00048\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0006\u00104R\u001c\u0010\b\u001a\n\u0012\u0004\u0012\u00020\u0007\u0018\u00010\u00048\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\b\u00104R\u0016\u0010\n\u001a\u0004\u0018\u00010\t8\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\n\u00105R\u0016\u0010\f\u001a\u0004\u0018\u00010\u000b8\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\f\u00106R\u0016\u0010\u000e\u001a\u0004\u0018\u00010\r8\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u000e\u00107R\u0016\u0010\u0010\u001a\u0004\u0018\u00010\u000f8\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0010\u00108R\u001c\u0010\u0012\u001a\n\u0012\u0004\u0012\u00020\u0011\u0018\u00010\u00048\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0012\u00104¨\u0006;"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageContents;", "Lcom/bendb/thrifty/a;", "", "message_text", "", "Lcom/x/dmv2/thriftjava/RichTextEntity;", "entities", "Lcom/x/dmv2/thriftjava/MessageAttachment;", "attachments", "Lcom/x/dmv2/thriftjava/ReplyingToPreview;", "replying_to_preview", "Lcom/x/dmv2/thriftjava/ForwardedMessage;", "forwarded_message", "Lcom/x/dmv2/thriftjava/SentFromSurface;", "sent_from", "Lcom/x/dmv2/thriftjava/QuickReply;", "quick_reply", "Lcom/x/dmv2/thriftjava/CallToAction;", "ctas", "<init>", "(Ljava/lang/String;Ljava/util/List;Ljava/util/List;Lcom/x/dmv2/thriftjava/ReplyingToPreview;Lcom/x/dmv2/thriftjava/ForwardedMessage;Lcom/x/dmv2/thriftjava/SentFromSurface;Lcom/x/dmv2/thriftjava/QuickReply;Ljava/util/List;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/String;", "component2", "()Ljava/util/List;", "component3", "component4", "()Lcom/x/dmv2/thriftjava/ReplyingToPreview;", "component5", "()Lcom/x/dmv2/thriftjava/ForwardedMessage;", "component6", "()Lcom/x/dmv2/thriftjava/SentFromSurface;", "component7", "()Lcom/x/dmv2/thriftjava/QuickReply;", "component8", "copy", "(Ljava/lang/String;Ljava/util/List;Ljava/util/List;Lcom/x/dmv2/thriftjava/ReplyingToPreview;Lcom/x/dmv2/thriftjava/ForwardedMessage;Lcom/x/dmv2/thriftjava/SentFromSurface;Lcom/x/dmv2/thriftjava/QuickReply;Ljava/util/List;)Lcom/x/dmv2/thriftjava/MessageContents;", "toString", "", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/String;", "Ljava/util/List;", "Lcom/x/dmv2/thriftjava/ReplyingToPreview;", "Lcom/x/dmv2/thriftjava/ForwardedMessage;", "Lcom/x/dmv2/thriftjava/SentFromSurface;", "Lcom/x/dmv2/thriftjava/QuickReply;", "Companion", "MessageContentsAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class MessageContents implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final List attachments;

    @JvmField
    @InterfaceC88465b
    public final List ctas;

    @JvmField
    @InterfaceC88465b
    public final List entities;

    @JvmField
    @InterfaceC88465b
    public final ForwardedMessage forwarded_message;

    @JvmField
    @InterfaceC88465b
    public final String message_text;

    @JvmField
    @InterfaceC88465b
    public final QuickReply quick_reply;

    @JvmField
    @InterfaceC88465b
    public final ReplyingToPreview replying_to_preview;

    @JvmField
    @InterfaceC88465b
    public final SentFromSurface sent_from;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new MessageContentsAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageContents$MessageContentsAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/MessageContents;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/MessageContents;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/MessageContents;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class MessageContentsAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public MessageContents m85653read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            String string = null;
            ArrayList arrayList = null;
            ArrayList arrayList2 = null;
            ReplyingToPreview replyingToPreview = null;
            ForwardedMessage forwardedMessage = null;
            SentFromSurface sentFromSurface = null;
            QuickReply quickReply = null;
            ArrayList arrayList3 = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b != 0) {
                    int i = 0;
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
                            if (b != 15) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                int i2 = protocol.mo14130a2().f38395b;
                                ArrayList arrayList4 = new ArrayList(i2);
                                while (i < i2) {
                                    arrayList4.add((RichTextEntity) RichTextEntity.ADAPTER.read(protocol));
                                    i++;
                                }
                                arrayList = arrayList4;
                                break;
                            }
                        case 3:
                            if (b != 15) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                int i3 = protocol.mo14130a2().f38395b;
                                ArrayList arrayList5 = new ArrayList(i3);
                                while (i < i3) {
                                    arrayList5.add((MessageAttachment) MessageAttachment.ADAPTER.read(protocol));
                                    i++;
                                }
                                arrayList2 = arrayList5;
                                break;
                            }
                        case 4:
                            if (b != 12) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                replyingToPreview = (ReplyingToPreview) ReplyingToPreview.ADAPTER.read(protocol);
                                break;
                            }
                        case 5:
                        default:
                            C11272a.m14141a(protocol, b);
                            break;
                        case 6:
                            if (b != 12) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                forwardedMessage = (ForwardedMessage) ForwardedMessage.ADAPTER.read(protocol);
                                break;
                            }
                        case 7:
                            if (b != 8) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                int iMo14132c4 = protocol.mo14132c4();
                                SentFromSurface sentFromSurfaceFindByValue = SentFromSurface.INSTANCE.findByValue(iMo14132c4);
                                if (sentFromSurfaceFindByValue == null) {
                                    throw new ThriftException(ThriftException.EnumC11260b.PROTOCOL_ERROR, C0031b.m45c(iMo14132c4, "Unexpected value for enum type SentFromSurface: "));
                                }
                                sentFromSurface = sentFromSurfaceFindByValue;
                                break;
                            }
                        case 8:
                            if (b != 12) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                quickReply = (QuickReply) QuickReply.ADAPTER.read(protocol);
                                break;
                            }
                        case 9:
                            if (b != 15) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                int i4 = protocol.mo14130a2().f38395b;
                                ArrayList arrayList6 = new ArrayList(i4);
                                while (i < i4) {
                                    arrayList6.add((CallToAction) CallToAction.ADAPTER.read(protocol));
                                    i++;
                                }
                                arrayList3 = arrayList6;
                                break;
                            }
                    }
                } else {
                    return new MessageContents(string, arrayList, arrayList2, replyingToPreview, forwardedMessage, sentFromSurface, quickReply, arrayList3);
                }
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a MessageContents struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("MessageContents");
            if (struct.message_text != null) {
                protocol.mo14136v3("message_text", 1, (byte) 11);
                protocol.mo14137w0(struct.message_text);
            }
            if (struct.entities != null) {
                protocol.mo14136v3("entities", 2, (byte) 15);
                protocol.mo14128X0((byte) 12, struct.entities.size());
                Iterator it = struct.entities.iterator();
                while (it.hasNext()) {
                    RichTextEntity.ADAPTER.write(protocol, (RichTextEntity) it.next());
                }
            }
            if (struct.attachments != null) {
                protocol.mo14136v3("attachments", 3, (byte) 15);
                protocol.mo14128X0((byte) 12, struct.attachments.size());
                Iterator it2 = struct.attachments.iterator();
                while (it2.hasNext()) {
                    MessageAttachment.ADAPTER.write(protocol, (MessageAttachment) it2.next());
                }
            }
            if (struct.replying_to_preview != null) {
                protocol.mo14136v3("replying_to_preview", 4, (byte) 12);
                ReplyingToPreview.ADAPTER.write(protocol, struct.replying_to_preview);
            }
            if (struct.forwarded_message != null) {
                protocol.mo14136v3("forwarded_message", 6, (byte) 12);
                ForwardedMessage.ADAPTER.write(protocol, struct.forwarded_message);
            }
            if (struct.sent_from != null) {
                protocol.mo14136v3("sent_from", 7, (byte) 8);
                protocol.mo14122C2(struct.sent_from.value);
            }
            if (struct.quick_reply != null) {
                protocol.mo14136v3("quick_reply", 8, (byte) 12);
                QuickReply.ADAPTER.write(protocol, struct.quick_reply);
            }
            if (struct.ctas != null) {
                protocol.mo14136v3("ctas", 9, (byte) 15);
                protocol.mo14128X0((byte) 12, struct.ctas.size());
                Iterator it3 = struct.ctas.iterator();
                while (it3.hasNext()) {
                    CallToAction.ADAPTER.write(protocol, (CallToAction) it3.next());
                }
            }
            protocol.mo14134i0();
        }
    }

    public MessageContents(@InterfaceC88465b String str, @InterfaceC88465b List list, @InterfaceC88465b List list2, @InterfaceC88465b ReplyingToPreview replyingToPreview, @InterfaceC88465b ForwardedMessage forwardedMessage, @InterfaceC88465b SentFromSurface sentFromSurface, @InterfaceC88465b QuickReply quickReply, @InterfaceC88465b List list3) {
        this.message_text = str;
        this.entities = list;
        this.attachments = list2;
        this.replying_to_preview = replyingToPreview;
        this.forwarded_message = forwardedMessage;
        this.sent_from = sentFromSurface;
        this.quick_reply = quickReply;
        this.ctas = list3;
    }

    public static /* synthetic */ MessageContents copy$default(MessageContents messageContents, String str, List list, List list2, ReplyingToPreview replyingToPreview, ForwardedMessage forwardedMessage, SentFromSurface sentFromSurface, QuickReply quickReply, List list3, int i, Object obj) {
        return messageContents.copy((i & 1) != 0 ? messageContents.message_text : str, (i & 2) != 0 ? messageContents.entities : list, (i & 4) != 0 ? messageContents.attachments : list2, (i & 8) != 0 ? messageContents.replying_to_preview : replyingToPreview, (i & 16) != 0 ? messageContents.forwarded_message : forwardedMessage, (i & 32) != 0 ? messageContents.sent_from : sentFromSurface, (i & 64) != 0 ? messageContents.quick_reply : quickReply, (i & 128) != 0 ? messageContents.ctas : list3);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final String getMessage_text() {
        return this.message_text;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final List getEntities() {
        return this.entities;
    }

    @InterfaceC88465b
    /* renamed from: component3, reason: from getter */
    public final List getAttachments() {
        return this.attachments;
    }

    @InterfaceC88465b
    /* renamed from: component4, reason: from getter */
    public final ReplyingToPreview getReplying_to_preview() {
        return this.replying_to_preview;
    }

    @InterfaceC88465b
    /* renamed from: component5, reason: from getter */
    public final ForwardedMessage getForwarded_message() {
        return this.forwarded_message;
    }

    @InterfaceC88465b
    /* renamed from: component6, reason: from getter */
    public final SentFromSurface getSent_from() {
        return this.sent_from;
    }

    @InterfaceC88465b
    /* renamed from: component7, reason: from getter */
    public final QuickReply getQuick_reply() {
        return this.quick_reply;
    }

    @InterfaceC88465b
    /* renamed from: component8, reason: from getter */
    public final List getCtas() {
        return this.ctas;
    }

    @InterfaceC88464a
    public final MessageContents copy(@InterfaceC88465b String message_text, @InterfaceC88465b List entities, @InterfaceC88465b List attachments, @InterfaceC88465b ReplyingToPreview replying_to_preview, @InterfaceC88465b ForwardedMessage forwarded_message, @InterfaceC88465b SentFromSurface sent_from, @InterfaceC88465b QuickReply quick_reply, @InterfaceC88465b List ctas) {
        return new MessageContents(message_text, entities, attachments, replying_to_preview, forwarded_message, sent_from, quick_reply, ctas);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof MessageContents)) {
            return false;
        }
        MessageContents messageContents = (MessageContents) other;
        return Intrinsics.m65267c(this.message_text, messageContents.message_text) && Intrinsics.m65267c(this.entities, messageContents.entities) && Intrinsics.m65267c(this.attachments, messageContents.attachments) && Intrinsics.m65267c(this.replying_to_preview, messageContents.replying_to_preview) && Intrinsics.m65267c(this.forwarded_message, messageContents.forwarded_message) && this.sent_from == messageContents.sent_from && Intrinsics.m65267c(this.quick_reply, messageContents.quick_reply) && Intrinsics.m65267c(this.ctas, messageContents.ctas);
    }

    public int hashCode() {
        String str = this.message_text;
        int iHashCode = (str == null ? 0 : str.hashCode()) * 31;
        List list = this.entities;
        int iHashCode2 = (iHashCode + (list == null ? 0 : list.hashCode())) * 31;
        List list2 = this.attachments;
        int iHashCode3 = (iHashCode2 + (list2 == null ? 0 : list2.hashCode())) * 31;
        ReplyingToPreview replyingToPreview = this.replying_to_preview;
        int iHashCode4 = (iHashCode3 + (replyingToPreview == null ? 0 : replyingToPreview.hashCode())) * 31;
        ForwardedMessage forwardedMessage = this.forwarded_message;
        int iHashCode5 = (iHashCode4 + (forwardedMessage == null ? 0 : forwardedMessage.hashCode())) * 31;
        SentFromSurface sentFromSurface = this.sent_from;
        int iHashCode6 = (iHashCode5 + (sentFromSurface == null ? 0 : sentFromSurface.hashCode())) * 31;
        QuickReply quickReply = this.quick_reply;
        int iHashCode7 = (iHashCode6 + (quickReply == null ? 0 : quickReply.hashCode())) * 31;
        List list3 = this.ctas;
        return iHashCode7 + (list3 != null ? list3.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        String str = this.message_text;
        List list = this.entities;
        List list2 = this.attachments;
        ReplyingToPreview replyingToPreview = this.replying_to_preview;
        ForwardedMessage forwardedMessage = this.forwarded_message;
        SentFromSurface sentFromSurface = this.sent_from;
        QuickReply quickReply = this.quick_reply;
        List list3 = this.ctas;
        StringBuilder sbM56620a = C68546u.m56620a("MessageContents(message_text=", str, ", entities=", list, ", attachments=");
        sbM56620a.append(list2);
        sbM56620a.append(", replying_to_preview=");
        sbM56620a.append(replyingToPreview);
        sbM56620a.append(", forwarded_message=");
        sbM56620a.append(forwardedMessage);
        sbM56620a.append(", sent_from=");
        sbM56620a.append(sentFromSurface);
        sbM56620a.append(", quick_reply=");
        sbM56620a.append(quickReply);
        sbM56620a.append(", ctas=");
        sbM56620a.append(list3);
        sbM56620a.append(Separators.RPAREN);
        return sbM56620a.toString();
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}
