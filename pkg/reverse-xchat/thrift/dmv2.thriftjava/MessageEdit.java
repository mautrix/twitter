package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import android.gov.nist.core.net.C0009a;
import androidx.camera.core.imagecapture.C0978g;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import java.io.IOException;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;
import kotlin.Metadata;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.Intrinsics;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u0000@\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\u000e\n\u0002\b\u0002\n\u0002\u0010 \n\u0002\u0018\u0002\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\n\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\u0007\b\u0086\b\u0018\u0000 !2\u00020\u0001:\u0002\"!B+\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0004\u001a\u0004\u0018\u00010\u0002\u0012\u000e\u0010\u0007\u001a\n\u0012\u0004\u0012\u00020\u0006\u0018\u00010\u0005¢\u0006\u0004\b\b\u0010\tJ\u0017\u0010\r\u001a\u00020\f2\u0006\u0010\u000b\u001a\u00020\nH\u0016¢\u0006\u0004\b\r\u0010\u000eJ\u0012\u0010\u000f\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u000f\u0010\u0010J\u0012\u0010\u0011\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\u0011\u0010\u0010J\u0018\u0010\u0012\u001a\n\u0012\u0004\u0012\u00020\u0006\u0018\u00010\u0005HÆ\u0003¢\u0006\u0004\b\u0012\u0010\u0013J:\u0010\u0014\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0004\u001a\u0004\u0018\u00010\u00022\u0010\b\u0002\u0010\u0007\u001a\n\u0012\u0004\u0012\u00020\u0006\u0018\u00010\u0005HÆ\u0001¢\u0006\u0004\b\u0014\u0010\u0015J\u0010\u0010\u0016\u001a\u00020\u0002HÖ\u0001¢\u0006\u0004\b\u0016\u0010\u0010J\u0010\u0010\u0018\u001a\u00020\u0017HÖ\u0001¢\u0006\u0004\b\u0018\u0010\u0019J\u001a\u0010\u001d\u001a\u00020\u001c2\b\u0010\u001b\u001a\u0004\u0018\u00010\u001aHÖ\u0003¢\u0006\u0004\b\u001d\u0010\u001eR\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010\u001fR\u0016\u0010\u0004\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0004\u0010\u001fR\u001c\u0010\u0007\u001a\n\u0012\u0004\u0012\u00020\u0006\u0018\u00010\u00058\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0007\u0010 ¨\u0006#"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEdit;", "Lcom/bendb/thrifty/a;", "", "message_sequence_id", "updated_text", "", "Lcom/x/dmv2/thriftjava/RichTextEntity;", "entities", "<init>", "(Ljava/lang/String;Ljava/lang/String;Ljava/util/List;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/String;", "component2", "component3", "()Ljava/util/List;", "copy", "(Ljava/lang/String;Ljava/lang/String;Ljava/util/List;)Lcom/x/dmv2/thriftjava/MessageEdit;", "toString", "", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/String;", "Ljava/util/List;", "Companion", "MessageEditAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class MessageEdit implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final List entities;

    @JvmField
    @InterfaceC88465b
    public final String message_sequence_id;

    @JvmField
    @InterfaceC88465b
    public final String updated_text;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new MessageEditAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MessageEdit$MessageEditAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/MessageEdit;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/MessageEdit;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/MessageEdit;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class MessageEditAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public MessageEdit m85655read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            String string = null;
            String string2 = null;
            ArrayList arrayList = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new MessageEdit(string, string2, arrayList);
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        if (s != 3) {
                            C11272a.m14141a(protocol, b);
                        } else if (b == 15) {
                            int i = protocol.mo14130a2().f38395b;
                            ArrayList arrayList2 = new ArrayList(i);
                            for (int i2 = 0; i2 < i; i2++) {
                                arrayList2.add((RichTextEntity) RichTextEntity.ADAPTER.read(protocol));
                            }
                            arrayList = arrayList2;
                        } else {
                            C11272a.m14141a(protocol, b);
                        }
                    } else if (b == 11) {
                        string2 = protocol.readString();
                    } else {
                        C11272a.m14141a(protocol, b);
                    }
                } else if (b == 11) {
                    string = protocol.readString();
                } else {
                    C11272a.m14141a(protocol, b);
                }
            }
        }

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a MessageEdit struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("MessageEdit");
            if (struct.message_sequence_id != null) {
                protocol.mo14136v3("message_sequence_id", 1, (byte) 11);
                protocol.mo14137w0(struct.message_sequence_id);
            }
            if (struct.updated_text != null) {
                protocol.mo14136v3("updated_text", 2, (byte) 11);
                protocol.mo14137w0(struct.updated_text);
            }
            if (struct.entities != null) {
                protocol.mo14136v3("entities", 3, (byte) 15);
                protocol.mo14128X0((byte) 12, struct.entities.size());
                Iterator it = struct.entities.iterator();
                while (it.hasNext()) {
                    RichTextEntity.ADAPTER.write(protocol, (RichTextEntity) it.next());
                }
            }
            protocol.mo14134i0();
        }
    }

    public MessageEdit(@InterfaceC88465b String str, @InterfaceC88465b String str2, @InterfaceC88465b List list) {
        this.message_sequence_id = str;
        this.updated_text = str2;
        this.entities = list;
    }

    public static /* synthetic */ MessageEdit copy$default(MessageEdit messageEdit, String str, String str2, List list, int i, Object obj) {
        if ((i & 1) != 0) {
            str = messageEdit.message_sequence_id;
        }
        if ((i & 2) != 0) {
            str2 = messageEdit.updated_text;
        }
        if ((i & 4) != 0) {
            list = messageEdit.entities;
        }
        return messageEdit.copy(str, str2, list);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final String getMessage_sequence_id() {
        return this.message_sequence_id;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final String getUpdated_text() {
        return this.updated_text;
    }

    @InterfaceC88465b
    /* renamed from: component3, reason: from getter */
    public final List getEntities() {
        return this.entities;
    }

    @InterfaceC88464a
    public final MessageEdit copy(@InterfaceC88465b String message_sequence_id, @InterfaceC88465b String updated_text, @InterfaceC88465b List entities) {
        return new MessageEdit(message_sequence_id, updated_text, entities);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof MessageEdit)) {
            return false;
        }
        MessageEdit messageEdit = (MessageEdit) other;
        return Intrinsics.m65267c(this.message_sequence_id, messageEdit.message_sequence_id) && Intrinsics.m65267c(this.updated_text, messageEdit.updated_text) && Intrinsics.m65267c(this.entities, messageEdit.entities);
    }

    public int hashCode() {
        String str = this.message_sequence_id;
        int iHashCode = (str == null ? 0 : str.hashCode()) * 31;
        String str2 = this.updated_text;
        int iHashCode2 = (iHashCode + (str2 == null ? 0 : str2.hashCode())) * 31;
        List list = this.entities;
        return iHashCode2 + (list != null ? list.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        String str = this.message_sequence_id;
        String str2 = this.updated_text;
        return C0978g.m1392a(C0009a.m11b("MessageEdit(message_sequence_id=", str, ", updated_text=", str2, ", entities="), this.entities, Separators.RPAREN);
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}
