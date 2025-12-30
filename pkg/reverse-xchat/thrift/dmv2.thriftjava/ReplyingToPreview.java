package com.x.dmv2.thriftjava;

import android.gov.nist.core.C0003b;
import android.gov.nist.core.Separators;
import android.gov.nist.javax.sip.clientauthutils.C0026b;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import com.chuckerteam.chucker.internal.data.har.log.entry.C11543d;
import java.io.IOException;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;
import kotlin.Metadata;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u0000J\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\t\n\u0000\n\u0002\u0010\u000e\n\u0000\n\u0002\u0010 \n\u0002\u0018\u0002\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0006\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\u000f\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\b\b\u0086\b\u0018\u0000 -2\u00020\u0001:\u0002.-BY\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0005\u001a\u0004\u0018\u00010\u0004\u0012\u000e\u0010\b\u001a\n\u0012\u0004\u0012\u00020\u0007\u0018\u00010\u0006\u0012\u000e\u0010\n\u001a\n\u0012\u0004\u0012\u00020\t\u0018\u00010\u0006\u0012\b\u0010\u000b\u001a\u0004\u0018\u00010\u0004\u0012\b\u0010\f\u001a\u0004\u0018\u00010\u0004\u0012\b\u0010\r\u001a\u0004\u0018\u00010\u0004¢\u0006\u0004\b\u000e\u0010\u000fJ\u0017\u0010\u0013\u001a\u00020\u00122\u0006\u0010\u0011\u001a\u00020\u0010H\u0016¢\u0006\u0004\b\u0013\u0010\u0014J\u0012\u0010\u0015\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0015\u0010\u0016J\u0012\u0010\u0017\u001a\u0004\u0018\u00010\u0004HÆ\u0003¢\u0006\u0004\b\u0017\u0010\u0018J\u0018\u0010\u0019\u001a\n\u0012\u0004\u0012\u00020\u0007\u0018\u00010\u0006HÆ\u0003¢\u0006\u0004\b\u0019\u0010\u001aJ\u0018\u0010\u001b\u001a\n\u0012\u0004\u0012\u00020\t\u0018\u00010\u0006HÆ\u0003¢\u0006\u0004\b\u001b\u0010\u001aJ\u0012\u0010\u001c\u001a\u0004\u0018\u00010\u0004HÆ\u0003¢\u0006\u0004\b\u001c\u0010\u0018J\u0012\u0010\u001d\u001a\u0004\u0018\u00010\u0004HÆ\u0003¢\u0006\u0004\b\u001d\u0010\u0018J\u0012\u0010\u001e\u001a\u0004\u0018\u00010\u0004HÆ\u0003¢\u0006\u0004\b\u001e\u0010\u0018Jp\u0010\u001f\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0005\u001a\u0004\u0018\u00010\u00042\u0010\b\u0002\u0010\b\u001a\n\u0012\u0004\u0012\u00020\u0007\u0018\u00010\u00062\u0010\b\u0002\u0010\n\u001a\n\u0012\u0004\u0012\u00020\t\u0018\u00010\u00062\n\b\u0002\u0010\u000b\u001a\u0004\u0018\u00010\u00042\n\b\u0002\u0010\f\u001a\u0004\u0018\u00010\u00042\n\b\u0002\u0010\r\u001a\u0004\u0018\u00010\u0004HÆ\u0001¢\u0006\u0004\b\u001f\u0010 J\u0010\u0010!\u001a\u00020\u0004HÖ\u0001¢\u0006\u0004\b!\u0010\u0018J\u0010\u0010#\u001a\u00020\"HÖ\u0001¢\u0006\u0004\b#\u0010$J\u001a\u0010(\u001a\u00020'2\b\u0010&\u001a\u0004\u0018\u00010%HÖ\u0003¢\u0006\u0004\b(\u0010)R\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010*R\u0016\u0010\u0005\u001a\u0004\u0018\u00010\u00048\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0005\u0010+R\u001c\u0010\b\u001a\n\u0012\u0004\u0012\u00020\u0007\u0018\u00010\u00068\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\b\u0010,R\u001c\u0010\n\u001a\n\u0012\u0004\u0012\u00020\t\u0018\u00010\u00068\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\n\u0010,R\u0016\u0010\u000b\u001a\u0004\u0018\u00010\u00048\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u000b\u0010+R\u0016\u0010\f\u001a\u0004\u0018\u00010\u00048\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\f\u0010+R\u0016\u0010\r\u001a\u0004\u0018\u00010\u00048\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\r\u0010+¨\u0006/"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/ReplyingToPreview;", "Lcom/bendb/thrifty/a;", "", "sender_id", "", "message_text", "", "Lcom/x/dmv2/thriftjava/RichTextEntity;", "entities", "Lcom/x/dmv2/thriftjava/MessageAttachment;", "attachments", "sender_display_name", "replying_to_message_sequence_id", "replying_to_message_id", "<init>", "(Ljava/lang/Long;Ljava/lang/String;Ljava/util/List;Ljava/util/List;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/Long;", "component2", "()Ljava/lang/String;", "component3", "()Ljava/util/List;", "component4", "component5", "component6", "component7", "copy", "(Ljava/lang/Long;Ljava/lang/String;Ljava/util/List;Ljava/util/List;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)Lcom/x/dmv2/thriftjava/ReplyingToPreview;", "toString", "", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/Long;", "Ljava/lang/String;", "Ljava/util/List;", "Companion", "ReplyingToPreviewAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class ReplyingToPreview implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final List attachments;

    @JvmField
    @InterfaceC88465b
    public final List entities;

    @JvmField
    @InterfaceC88465b
    public final String message_text;

    @JvmField
    @InterfaceC88465b
    public final String replying_to_message_id;

    @JvmField
    @InterfaceC88465b
    public final String replying_to_message_sequence_id;

    @JvmField
    @InterfaceC88465b
    public final String sender_display_name;

    @JvmField
    @InterfaceC88465b
    public final Long sender_id;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new ReplyingToPreviewAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/ReplyingToPreview$ReplyingToPreviewAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/ReplyingToPreview;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/ReplyingToPreview;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/ReplyingToPreview;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class ReplyingToPreviewAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public ReplyingToPreview m85669read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Long lValueOf = null;
            String string = null;
            ArrayList arrayList = null;
            ArrayList arrayList2 = null;
            String string2 = null;
            String string3 = null;
            String string4 = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b != 0) {
                    int i = 0;
                    switch (c11265cMo14127V2.f38393b) {
                        case 1:
                            if (b != 10) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                lValueOf = Long.valueOf(protocol.mo14124H0());
                                break;
                            }
                        case 2:
                            if (b != 11) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                string = protocol.readString();
                                break;
                            }
                        case 3:
                            if (b != 15) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                int i2 = protocol.mo14130a2().f38395b;
                                ArrayList arrayList3 = new ArrayList(i2);
                                while (i < i2) {
                                    arrayList3.add((RichTextEntity) RichTextEntity.ADAPTER.read(protocol));
                                    i++;
                                }
                                arrayList = arrayList3;
                                break;
                            }
                        case 4:
                            if (b != 15) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                int i3 = protocol.mo14130a2().f38395b;
                                ArrayList arrayList4 = new ArrayList(i3);
                                while (i < i3) {
                                    arrayList4.add((MessageAttachment) MessageAttachment.ADAPTER.read(protocol));
                                    i++;
                                }
                                arrayList2 = arrayList4;
                                break;
                            }
                        case 5:
                            if (b != 11) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                string2 = protocol.readString();
                                break;
                            }
                        case 6:
                            if (b != 11) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                string3 = protocol.readString();
                                break;
                            }
                        case 7:
                            if (b != 11) {
                                C11272a.m14141a(protocol, b);
                                break;
                            } else {
                                string4 = protocol.readString();
                                break;
                            }
                        default:
                            C11272a.m14141a(protocol, b);
                            break;
                    }
                } else {
                    return new ReplyingToPreview(lValueOf, string, arrayList, arrayList2, string2, string3, string4);
                }
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a ReplyingToPreview struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("ReplyingToPreview");
            if (struct.sender_id != null) {
                protocol.mo14136v3("sender_id", 1, (byte) 10);
                protocol.mo14121B3(struct.sender_id.longValue());
            }
            if (struct.message_text != null) {
                protocol.mo14136v3("message_text", 2, (byte) 11);
                protocol.mo14137w0(struct.message_text);
            }
            if (struct.entities != null) {
                protocol.mo14136v3("entities", 3, (byte) 15);
                protocol.mo14128X0((byte) 12, struct.entities.size());
                Iterator it = struct.entities.iterator();
                while (it.hasNext()) {
                    RichTextEntity.ADAPTER.write(protocol, (RichTextEntity) it.next());
                }
            }
            if (struct.attachments != null) {
                protocol.mo14136v3("attachments", 4, (byte) 15);
                protocol.mo14128X0((byte) 12, struct.attachments.size());
                Iterator it2 = struct.attachments.iterator();
                while (it2.hasNext()) {
                    MessageAttachment.ADAPTER.write(protocol, (MessageAttachment) it2.next());
                }
            }
            if (struct.sender_display_name != null) {
                protocol.mo14136v3("sender_display_name", 5, (byte) 11);
                protocol.mo14137w0(struct.sender_display_name);
            }
            if (struct.replying_to_message_sequence_id != null) {
                protocol.mo14136v3("replying_to_message_sequence_id", 6, (byte) 11);
                protocol.mo14137w0(struct.replying_to_message_sequence_id);
            }
            if (struct.replying_to_message_id != null) {
                protocol.mo14136v3("replying_to_message_id", 7, (byte) 11);
                protocol.mo14137w0(struct.replying_to_message_id);
            }
            protocol.mo14134i0();
        }
    }

    public ReplyingToPreview(@InterfaceC88465b Long l, @InterfaceC88465b String str, @InterfaceC88465b List list, @InterfaceC88465b List list2, @InterfaceC88465b String str2, @InterfaceC88465b String str3, @InterfaceC88465b String str4) {
        this.sender_id = l;
        this.message_text = str;
        this.entities = list;
        this.attachments = list2;
        this.sender_display_name = str2;
        this.replying_to_message_sequence_id = str3;
        this.replying_to_message_id = str4;
    }

    public static /* synthetic */ ReplyingToPreview copy$default(ReplyingToPreview replyingToPreview, Long l, String str, List list, List list2, String str2, String str3, String str4, int i, Object obj) {
        if ((i & 1) != 0) {
            l = replyingToPreview.sender_id;
        }
        if ((i & 2) != 0) {
            str = replyingToPreview.message_text;
        }
        String str5 = str;
        if ((i & 4) != 0) {
            list = replyingToPreview.entities;
        }
        List list3 = list;
        if ((i & 8) != 0) {
            list2 = replyingToPreview.attachments;
        }
        List list4 = list2;
        if ((i & 16) != 0) {
            str2 = replyingToPreview.sender_display_name;
        }
        String str6 = str2;
        if ((i & 32) != 0) {
            str3 = replyingToPreview.replying_to_message_sequence_id;
        }
        String str7 = str3;
        if ((i & 64) != 0) {
            str4 = replyingToPreview.replying_to_message_id;
        }
        return replyingToPreview.copy(l, str5, list3, list4, str6, str7, str4);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final Long getSender_id() {
        return this.sender_id;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final String getMessage_text() {
        return this.message_text;
    }

    @InterfaceC88465b
    /* renamed from: component3, reason: from getter */
    public final List getEntities() {
        return this.entities;
    }

    @InterfaceC88465b
    /* renamed from: component4, reason: from getter */
    public final List getAttachments() {
        return this.attachments;
    }

    @InterfaceC88465b
    /* renamed from: component5, reason: from getter */
    public final String getSender_display_name() {
        return this.sender_display_name;
    }

    @InterfaceC88465b
    /* renamed from: component6, reason: from getter */
    public final String getReplying_to_message_sequence_id() {
        return this.replying_to_message_sequence_id;
    }

    @InterfaceC88465b
    /* renamed from: component7, reason: from getter */
    public final String getReplying_to_message_id() {
        return this.replying_to_message_id;
    }

    @InterfaceC88464a
    public final ReplyingToPreview copy(@InterfaceC88465b Long sender_id, @InterfaceC88465b String message_text, @InterfaceC88465b List entities, @InterfaceC88465b List attachments, @InterfaceC88465b String sender_display_name, @InterfaceC88465b String replying_to_message_sequence_id, @InterfaceC88465b String replying_to_message_id) {
        return new ReplyingToPreview(sender_id, message_text, entities, attachments, sender_display_name, replying_to_message_sequence_id, replying_to_message_id);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof ReplyingToPreview)) {
            return false;
        }
        ReplyingToPreview replyingToPreview = (ReplyingToPreview) other;
        return Intrinsics.m65267c(this.sender_id, replyingToPreview.sender_id) && Intrinsics.m65267c(this.message_text, replyingToPreview.message_text) && Intrinsics.m65267c(this.entities, replyingToPreview.entities) && Intrinsics.m65267c(this.attachments, replyingToPreview.attachments) && Intrinsics.m65267c(this.sender_display_name, replyingToPreview.sender_display_name) && Intrinsics.m65267c(this.replying_to_message_sequence_id, replyingToPreview.replying_to_message_sequence_id) && Intrinsics.m65267c(this.replying_to_message_id, replyingToPreview.replying_to_message_id);
    }

    public int hashCode() {
        Long l = this.sender_id;
        int iHashCode = (l == null ? 0 : l.hashCode()) * 31;
        String str = this.message_text;
        int iHashCode2 = (iHashCode + (str == null ? 0 : str.hashCode())) * 31;
        List list = this.entities;
        int iHashCode3 = (iHashCode2 + (list == null ? 0 : list.hashCode())) * 31;
        List list2 = this.attachments;
        int iHashCode4 = (iHashCode3 + (list2 == null ? 0 : list2.hashCode())) * 31;
        String str2 = this.sender_display_name;
        int iHashCode5 = (iHashCode4 + (str2 == null ? 0 : str2.hashCode())) * 31;
        String str3 = this.replying_to_message_sequence_id;
        int iHashCode6 = (iHashCode5 + (str3 == null ? 0 : str3.hashCode())) * 31;
        String str4 = this.replying_to_message_id;
        return iHashCode6 + (str4 != null ? str4.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        Long l = this.sender_id;
        String str = this.message_text;
        List list = this.entities;
        List list2 = this.attachments;
        String str2 = this.sender_display_name;
        String str3 = this.replying_to_message_sequence_id;
        String str4 = this.replying_to_message_id;
        StringBuilder sb = new StringBuilder("ReplyingToPreview(sender_id=");
        sb.append(l);
        sb.append(", message_text=");
        sb.append(str);
        sb.append(", entities=");
        C11543d.m14632b(sb, list, ", attachments=", list2, ", sender_display_name=");
        C0026b.m37b(sb, str2, ", replying_to_message_sequence_id=", str3, ", replying_to_message_id=");
        return C0003b.m4b(sb, str4, Separators.RPAREN);
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}
