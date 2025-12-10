package com.x.dmv2.thriftjava;

import android.gov.nist.core.Separators;
import com.bendb.thrifty.InterfaceC11261a;
import com.bendb.thrifty.kotlin.InterfaceC11262a;
import com.bendb.thrifty.protocol.C11265c;
import com.bendb.thrifty.protocol.InterfaceC11268f;
import com.bendb.thrifty.util.C11272a;
import java.io.IOException;
import kotlin.Metadata;
import kotlin.jvm.JvmField;
import kotlin.jvm.internal.Intrinsics;
import okio.C87081h;
import org.jetbrains.annotations.InterfaceC88464a;
import org.jetbrains.annotations.InterfaceC88465b;

@Metadata(m64929d1 = {"\u0000:\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0010\u000e\n\u0000\n\u0002\u0018\u0002\n\u0002\b\u0003\n\u0002\u0018\u0002\n\u0000\n\u0002\u0010\u0002\n\u0002\b\t\n\u0002\u0010\b\n\u0002\b\u0002\n\u0002\u0010\u0000\n\u0000\n\u0002\u0010\u000b\n\u0002\b\u0007\b\u0086\b\u0018\u0000 \u001e2\u00020\u0001:\u0002\u001f\u001eB\u001b\u0012\b\u0010\u0003\u001a\u0004\u0018\u00010\u0002\u0012\b\u0010\u0005\u001a\u0004\u0018\u00010\u0004¢\u0006\u0004\b\u0006\u0010\u0007J\u0017\u0010\u000b\u001a\u00020\n2\u0006\u0010\t\u001a\u00020\bH\u0016¢\u0006\u0004\b\u000b\u0010\fJ\u0012\u0010\r\u001a\u0004\u0018\u00010\u0002HÆ\u0003¢\u0006\u0004\b\r\u0010\u000eJ\u0012\u0010\u000f\u001a\u0004\u0018\u00010\u0004HÆ\u0003¢\u0006\u0004\b\u000f\u0010\u0010J(\u0010\u0011\u001a\u00020\u00002\n\b\u0002\u0010\u0003\u001a\u0004\u0018\u00010\u00022\n\b\u0002\u0010\u0005\u001a\u0004\u0018\u00010\u0004HÆ\u0001¢\u0006\u0004\b\u0011\u0010\u0012J\u0010\u0010\u0013\u001a\u00020\u0002HÖ\u0001¢\u0006\u0004\b\u0013\u0010\u000eJ\u0010\u0010\u0015\u001a\u00020\u0014HÖ\u0001¢\u0006\u0004\b\u0015\u0010\u0016J\u001a\u0010\u001a\u001a\u00020\u00192\b\u0010\u0018\u001a\u0004\u0018\u00010\u0017HÖ\u0003¢\u0006\u0004\b\u001a\u0010\u001bR\u0016\u0010\u0003\u001a\u0004\u0018\u00010\u00028\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0003\u0010\u001cR\u0016\u0010\u0005\u001a\u0004\u0018\u00010\u00048\u0006X\u0087\u0004¢\u0006\u0006\n\u0004\b\u0005\u0010\u001d¨\u0006 "}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MoneyAttachment;", "Lcom/bendb/thrifty/a;", "", "fallbackText", "Lokio/h;", "payload", "<init>", "(Ljava/lang/String;Lokio/h;)V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "", "write", "(Lcom/bendb/thrifty/protocol/f;)V", "component1", "()Ljava/lang/String;", "component2", "()Lokio/h;", "copy", "(Ljava/lang/String;Lokio/h;)Lcom/x/dmv2/thriftjava/MoneyAttachment;", "toString", "", "hashCode", "()I", "", "other", "", "equals", "(Ljava/lang/Object;)Z", "Ljava/lang/String;", "Lokio/h;", "Companion", "MoneyAttachmentAdapter", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
/* loaded from: classes4.dex */
public final /* data */ class MoneyAttachment implements InterfaceC11261a {

    @JvmField
    @InterfaceC88465b
    public final String fallbackText;

    @JvmField
    @InterfaceC88465b
    public final C87081h payload;

    @JvmField
    @InterfaceC88464a
    public static final InterfaceC11262a ADAPTER = new MoneyAttachmentAdapter();

    @Metadata(m64929d1 = {"\u0000 \n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\u0018\u0002\n\u0002\b\u0002\n\u0002\u0018\u0002\n\u0002\b\u0004\n\u0002\u0010\u0002\n\u0002\b\u0003\b\u0002\u0018\u00002\b\u0012\u0004\u0012\u00020\u00020\u0001B\u0007¢\u0006\u0004\b\u0003\u0010\u0004J\u0017\u0010\u0007\u001a\u00020\u00022\u0006\u0010\u0006\u001a\u00020\u0005H\u0016¢\u0006\u0004\b\u0007\u0010\bJ\u001f\u0010\u000b\u001a\u00020\n2\u0006\u0010\u0006\u001a\u00020\u00052\u0006\u0010\t\u001a\u00020\u0002H\u0016¢\u0006\u0004\b\u000b\u0010\f¨\u0006\r"}, m64930d2 = {"Lcom/x/dmv2/thriftjava/MoneyAttachment$MoneyAttachmentAdapter;", "Lcom/bendb/thrifty/kotlin/a;", "Lcom/x/dmv2/thriftjava/MoneyAttachment;", "<init>", "()V", "Lcom/bendb/thrifty/protocol/f;", "protocol", "read", "(Lcom/bendb/thrifty/protocol/f;)Lcom/x/dmv2/thriftjava/MoneyAttachment;", "struct", "", "write", "(Lcom/bendb/thrifty/protocol/f;Lcom/x/dmv2/thriftjava/MoneyAttachment;)V", "-subsystem-dm-thrift"}, m64931k = 1, m64932mv = {2, 1, 0}, m64934xi = 48)
    public static final class MoneyAttachmentAdapter implements InterfaceC11262a {
        @InterfaceC88464a
        /* renamed from: read, reason: merged with bridge method [inline-methods] */
        public MoneyAttachment m85950read(@InterfaceC88464a InterfaceC11268f protocol) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            String string = null;
            C87081h c87081hMo14123G1 = null;
            while (true) {
                C11265c c11265cMo14127V2 = protocol.mo14127V2();
                byte b = c11265cMo14127V2.f38392a;
                if (b == 0) {
                    return new MoneyAttachment(string, c87081hMo14123G1);
                }
                short s = c11265cMo14127V2.f38393b;
                if (s != 1) {
                    if (s != 2) {
                        C11272a.m14141a(protocol, b);
                    } else if (b == 11) {
                        c87081hMo14123G1 = protocol.mo14123G1();
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

        public void write(@InterfaceC88464a InterfaceC11268f protocol, @InterfaceC88464a MoneyAttachment struct) throws IOException {
            Intrinsics.m65272h(protocol, "protocol");
            Intrinsics.m65272h(struct, "struct");
            protocol.mo14129Y2("MoneyAttachment");
            if (struct.fallbackText != null) {
                protocol.mo14136v3("fallbackText", 1, (byte) 11);
                protocol.mo14137w0(struct.fallbackText);
            }
            if (struct.payload != null) {
                protocol.mo14136v3("payload", 2, (byte) 11);
                protocol.mo14138y0(struct.payload);
            }
            protocol.mo14134i0();
        }
    }

    public MoneyAttachment(@InterfaceC88465b String str, @InterfaceC88465b C87081h c87081h) {
        this.fallbackText = str;
        this.payload = c87081h;
    }

    public static /* synthetic */ MoneyAttachment copy$default(MoneyAttachment moneyAttachment, String str, C87081h c87081h, int i, Object obj) {
        if ((i & 1) != 0) {
            str = moneyAttachment.fallbackText;
        }
        if ((i & 2) != 0) {
            c87081h = moneyAttachment.payload;
        }
        return moneyAttachment.copy(str, c87081h);
    }

    @InterfaceC88465b
    /* renamed from: component1, reason: from getter */
    public final String getFallbackText() {
        return this.fallbackText;
    }

    @InterfaceC88465b
    /* renamed from: component2, reason: from getter */
    public final C87081h getPayload() {
        return this.payload;
    }

    @InterfaceC88464a
    public final MoneyAttachment copy(@InterfaceC88465b String fallbackText, @InterfaceC88465b C87081h payload) {
        return new MoneyAttachment(fallbackText, payload);
    }

    public boolean equals(@InterfaceC88465b Object other) {
        if (this == other) {
            return true;
        }
        if (!(other instanceof MoneyAttachment)) {
            return false;
        }
        MoneyAttachment moneyAttachment = (MoneyAttachment) other;
        return Intrinsics.m65267c(this.fallbackText, moneyAttachment.fallbackText) && Intrinsics.m65267c(this.payload, moneyAttachment.payload);
    }

    public int hashCode() {
        String str = this.fallbackText;
        int iHashCode = (str == null ? 0 : str.hashCode()) * 31;
        C87081h c87081h = this.payload;
        return iHashCode + (c87081h != null ? c87081h.hashCode() : 0);
    }

    @InterfaceC88464a
    public String toString() {
        return "MoneyAttachment(fallbackText=" + this.fallbackText + ", payload=" + this.payload + Separators.RPAREN;
    }

    public void write(@InterfaceC88464a InterfaceC11268f protocol) {
        Intrinsics.m65272h(protocol, "protocol");
        ADAPTER.write(protocol, this);
    }
}